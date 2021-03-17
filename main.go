/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	kafkav1beta1 "github.com/DoodleScheduling/k8skafka-controller/api/v1beta1"
	"github.com/DoodleScheduling/k8skafka-controller/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

var (
	metricsAddr             = ":9556"
	probesAddr              = ":9557"
	enableLeaderElection    = true
	leaderElectionNamespace = ""
	namespaces              = ""
	concurrent              = 4
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = kafkav1beta1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	flag.StringVar(&metricsAddr, "metrics-addr", ":9556", "The address of the metric endpoint binds to.")
	flag.StringVar(&probesAddr, "probe-addr", ":9557", "The address of the probe endpoints bind to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", true,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&leaderElectionNamespace, "leader-election-namespace", "",
		"Specify a different leader election namespace. It will use the one where the controller is deployed by default.")
	flag.StringVar(&namespaces, "namespaces", "",
		"The controller listens by default for all namespaces. This may be limited to a comma delimted list of dedicated namespaces.")
	flag.IntVar(&concurrent, "concurrent", 4,
		"The number of concurrent reconcile workers. By default this is 4.")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Import flags into viper and bind them to env vars
	// flags are converted to upper-case, - is replaced with _
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		setupLog.Error(err, "Failed parsing command line arguments")
		os.Exit(1)
	}

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	opts := ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      viper.GetString("metrics-addr"),
		HealthProbeBindAddress:  viper.GetString("probe-addr"),
		Port:                    9443,
		LeaderElection:          viper.GetBool("enable-leader-election"),
		LeaderElectionNamespace: viper.GetString("leader-election-namespace"),
		LeaderElectionID:        "1e457812.doodle.com",
	}

	ns := strings.Split(viper.GetString("namespaces"), ",")
	if len(ns) > 0 && ns[0] != "" {
		opts.NewCache = cache.MultiNamespacedCacheBuilder(ns)
		setupLog.Info("watching dedicated namespaces", "namespaces", ns)
	} else {
		setupLog.Info("watching all namespaces")
	}
	if leaderElectionNamespace != "" {
		opts.LeaderElectionNamespace = leaderElectionNamespace
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), opts)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Add liveness probe
	err = mgr.AddHealthzCheck("healthz", healthz.Ping)
	if err != nil {
		setupLog.Error(err, "Could not add liveness probe")
		os.Exit(1)
	}

	// Add readiness probe
	err = mgr.AddReadyzCheck("readyz", healthz.Ping)
	if err != nil {
		setupLog.Error(err, "Could not add readiness probe")
		os.Exit(1)
	}

	if err = (&controllers.KafkaTopicReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("KafkaTopic"),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("KafkaTopic"),
	}).SetupWithManager(mgr, controllers.KafkaTopicReconcilerOptions{MaxConcurrentReconciles: viper.GetInt("concurrent")}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KafkaTopic")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}