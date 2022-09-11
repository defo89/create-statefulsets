package cmd

import (
	"context"
	goflag "flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/defo89/create-statefulsets/pkg/config"
	"github.com/defo89/create-statefulsets/pkg/statefulset"
	"github.com/defo89/create-statefulsets/pkg/volumeclaim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var rootCmd = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  "",
	RunE:  runRootCmd,
}

var cfg = config.Config{}

func init() {
	viper.AutomaticEnv()
	rootCmd.PersistentFlags().StringVar(&cfg.KubeConfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	rootCmd.PersistentFlags().BoolVar(&cfg.Create, "create", true, "create statefulsets")
	rootCmd.PersistentFlags().BoolVar(&cfg.Delete, "delete", false, "delete statefulsets")
	rootCmd.PersistentFlags().StringVar(&cfg.StorageClass, "storage-class", "local-path", "storage class")
	rootCmd.PersistentFlags().StringVar(&cfg.PVCSize, "pvc-size", "5Mi", "PVC size")
	rootCmd.PersistentFlags().StringVar(&cfg.ImageName, "image", "nginx", "image name")
	rootCmd.PersistentFlags().StringVar(&cfg.ImageTag, "tag", "latest", "image tag")
	rootCmd.PersistentFlags().StringVar(&cfg.Namespace, "namespace", "default", "namespace")
	rootCmd.PersistentFlags().IntVar(&cfg.Count, "count", 1, "amount of statefulsets to create")
	rootCmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)
	_ = viper.BindPFlags(rootCmd.PersistentFlags())
}

func runRootCmd(cmd *cobra.Command, args []string) error {

	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = cfg.KubeConfig
	}
	kconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(kconfig)
	if err != nil {
		panic(err)
	}

	pvcClient := clientset.CoreV1().PersistentVolumeClaims(cfg.Namespace)
	stsClient := clientset.AppsV1().StatefulSets(cfg.Namespace)

	for i := 1; i <= cfg.Count; i++ {
		sts, err := statefulset.CreateStatefulsetObject(i, cfg.ImageName, cfg.ImageTag)
		if err != nil {
			panic(err)
		}

		pvc, err := volumeclaim.CreateStatefulsetObject(i, cfg.StorageClass, cfg.PVCSize)
		if err != nil {
			panic(err)
		}

		if cfg.Delete {
			deletePolicy := metav1.DeletePropagationForeground

			_, err := stsClient.Delete(context.TODO(), sts.Name, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy}), err
			if err != nil {
				panic(err)
			}

			_, err = pvcClient.Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy}), err
			if err != nil {
				panic(err)
			}

		} else {
			_, err := stsClient.Create(context.TODO(), sts, metav1.CreateOptions{})
			if err != nil {
				panic(err)
			}

			_, err = pvcClient.Create(context.TODO(), pvc, metav1.CreateOptions{})
			if err != nil {
				panic(err)
			}
		}
	}

	if cfg.Delete {
		fmt.Printf("Deleted %d statefulsets\n", cfg.Count)
	} else {
		fmt.Printf("Created %d statefulsets\n", cfg.Count)
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
