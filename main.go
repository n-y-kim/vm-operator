// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"crypto/tls"
	"flag"
	"net/http"
	"net/http/pprof"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/client-go/util/flowcontrol"
	klog "k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"

	"github.com/vmware-tanzu/vm-operator/api/v1alpha1"
	"github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	"github.com/vmware-tanzu/vm-operator/controllers"
	"github.com/vmware-tanzu/vm-operator/pkg"
	"github.com/vmware-tanzu/vm-operator/pkg/context"
	"github.com/vmware-tanzu/vm-operator/pkg/lib"
	"github.com/vmware-tanzu/vm-operator/pkg/manager"
	"github.com/vmware-tanzu/vm-operator/webhooks"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlmgr "sigs.k8s.io/controller-runtime/pkg/manager"
	ctrlsig "sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	// +kubebuilder:scaffold:imports
)

var (
	defaultProfilerAddr     = ":8073"
	defaultRateLimiterQPS   = 500
	defaultRateLimiterBurst = 1000

	defaultSyncPeriod                   = manager.DefaultSyncPeriod
	defaultMaxConcurrentReconciles      = manager.DefaultMaxConcurrentReconciles
	defaultLeaderElectionID             = manager.DefaultLeaderElectionID
	defaultPodNamespace                 = manager.DefaultPodNamespace
	defaultPodName                      = manager.DefaultPodName
	defaultPodServiceAccountName        = manager.DefaultPodServiceAccountName
	defaultWebhookServiceContainerPort  = manager.DefaultWebhookServiceContainerPort
	defaultWebhookServiceNamespace      = manager.DefaultWebhookServiceNamespace
	defaultWebhookServiceName           = manager.DefaultWebhookServiceName
	defaultWebhookSecretNamespace       = manager.DefaultWebhookSecretNamespace
	defaultWebhookSecretName            = manager.DefaultWebhookSecretName
	defaultWebhookSecretVolumeMountPath = manager.DefaultWebhookSecretVolumeMountPath
	defaultWatchNamespace               = manager.DefaultWatchNamespace
	defaultContainerNode                = manager.DefaultContainerNode
)

const (
	// serverKeyName is the name of the server private key.
	serverKeyName = "tls.key"
	// serverCertName is the name of the serving certificate.
	serverCertName = "tls.crt"
)

func init() {
	if v := os.Getenv("PROFILER_ADDR"); v != "" {
		defaultProfilerAddr = v
	}
	if v, err := strconv.Atoi(os.Getenv("RATE_LIMIT_QPS")); err == nil {
		defaultRateLimiterQPS = v
	}
	if v, err := strconv.Atoi(os.Getenv("RATE_LIMIT_BURST")); err == nil {
		defaultRateLimiterBurst = v
	}
	if v, err := time.ParseDuration(os.Getenv("SYNC_PERIOD")); err == nil {
		defaultSyncPeriod = v
	}
	if v, err := strconv.Atoi(os.Getenv("MAX_CONCURRENT_RECONCILES")); err == nil {
		defaultMaxConcurrentReconciles = v
	}
	if v := os.Getenv("LEADER_ELECTION_ID"); v != "" {
		defaultLeaderElectionID = v
	}
	if v := os.Getenv("POD_NAMESPACE"); v != "" {
		defaultPodNamespace = v
	}
	if v := os.Getenv("POD_NAME"); v != "" {
		defaultPodName = v
	}
	if v := os.Getenv("POD_SERVICE_ACCOUNT_NAME"); v != "" {
		defaultPodServiceAccountName = v
	}
	if v, err := strconv.Atoi(os.Getenv("WEBHOOK_SERVICE_CONTAINER_PORT")); err == nil {
		defaultWebhookServiceContainerPort = v
	}
	if v := os.Getenv("WEBHOOK_SERVICE_NAMESPACE"); v != "" {
		defaultWebhookServiceNamespace = v
	}
	if v := os.Getenv("WEBHOOK_SERVICE_NAME"); v != "" {
		defaultWebhookServiceName = v
	}
	if v := os.Getenv("WEBHOOK_SECRET_NAMESPACE"); v != "" {
		defaultWebhookSecretNamespace = v
	}
	if v := os.Getenv("WEBHOOK_SECRET_NAME"); v != "" {
		defaultWebhookSecretName = v
	}
	if v := os.Getenv("WATCH_NAMESPACE"); v != "" {
		defaultWatchNamespace = v
	}
	defaultContainerNode, _ = strconv.ParseBool(os.Getenv("CONTAINER_NODE"))
}

func main() {
	klog.InitFlags(nil)
	ctrllog.SetLogger(klogr.New())
	setupLog := ctrllog.Log.WithName("entrypoint")

	setupLog.Info("Starting VM Operator controller", "version", pkg.BuildVersion,
		"buildnumber", pkg.BuildNumber, "buildtype", pkg.BuildType, "commit", pkg.BuildCommit)

	profilerAddress := flag.String(
		"profiler-address",
		defaultProfilerAddr,
		"Bind address to expose the pprof profiler.",
	)
	rateLimiterQPS := flag.Int(
		"rate-limit-requests-per-second",
		defaultRateLimiterQPS,
		"The default number of requests per second to configure the k8s client rate limiter to allow.",
	)
	rateLimiterBurst := flag.Int(
		"rate-limit-max-requests",
		defaultRateLimiterBurst,
		"The default number of maximum burst requests per second to configure the k8s client rate limiter to allow.",
	)

	var managerOpts manager.Options

	flag.StringVar(
		&managerOpts.MetricsAddr,
		"metrics-addr",
		":8083",
		"The address the metric endpoint binds to.")
	flag.StringVar(
		&managerOpts.HealthProbeBindAddress,
		"health-addr",
		":9445",
		"The address the health probe endpoint binds to.")
	flag.BoolVar(
		&managerOpts.LeaderElectionEnabled,
		"enable-leader-election",
		true,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(
		&managerOpts.LeaderElectionID,
		"leader-election-id",
		defaultLeaderElectionID,
		"Name of the config map to use as the locking resource when configuring leader election.")
	flag.StringVar(
		&managerOpts.WatchNamespace,
		"watch-namespace",
		defaultWatchNamespace,
		"Namespace that the controller watches to reconcile vm operator objects. If unspecified, the controller watches for vm operator objects across all namespaces.")
	flag.DurationVar(
		&managerOpts.SyncPeriod,
		"sync-period",
		defaultSyncPeriod,
		"The interval at which objects are synchronized.")
	flag.IntVar(
		&managerOpts.MaxConcurrentReconciles,
		"max-concurrent-reconciles",
		defaultMaxConcurrentReconciles,
		"The maximum number of allowed, concurrent reconciles.")
	flag.StringVar(
		&managerOpts.PodNamespace,
		"pod-namespace",
		defaultPodNamespace,
		"The namespace in which the pod running the controller manager is located.")
	flag.StringVar(
		&managerOpts.PodName,
		"pod-name",
		defaultPodName,
		"The name of the pod running the controller manager.")
	flag.StringVar(
		&managerOpts.PodServiceAccountName,
		"pod-service-account-name",
		defaultPodServiceAccountName,
		"The service account name of the pod running the controller manager.")
	flag.IntVar(
		&managerOpts.WebhookServiceContainerPort,
		"webhook-service-container-port",
		defaultWebhookServiceContainerPort,
		"The the port on which the webhook service expects the webhook server to listen for incoming requests.")
	flag.StringVar(
		&managerOpts.WebhookServiceNamespace,
		"webhook-service-namespace",
		defaultWebhookServiceNamespace,
		"The namespace in which the webhook service is located.")
	flag.StringVar(
		&managerOpts.WebhookServiceName,
		"webhook-service-name",
		defaultWebhookServiceName,
		"The name of the webhook service.")
	flag.StringVar(
		&managerOpts.WebhookSecretNamespace,
		"webhook-secret-namespace",
		defaultWebhookSecretNamespace,
		"The namespace in which the webhook secret is located.")
	flag.StringVar(
		&managerOpts.WebhookSecretName,
		"webhook-secret-name",
		defaultWebhookSecretName,
		"The name of the webhook secret.")
	flag.StringVar(
		&managerOpts.WebhookSecretVolumeMountPath,
		"webhook-secret-volume-mount-path",
		defaultWebhookSecretVolumeMountPath,
		"The filesystem path to which the webhook secret is mounted.")
	flag.BoolVar(
		&managerOpts.ContainerNode,
		"container-node",
		defaultContainerNode,
		"Should be true if we're running nodes in containers (with vcsim).",
	)

	flag.Parse()

	if managerOpts.WatchNamespace != "" {
		setupLog.Info(
			"Watching objects only in namespace for reconciliation",
			"namespace", managerOpts.WatchNamespace)
	}

	if *rateLimiterQPS != 0 || *rateLimiterBurst != 0 {
		cfg := ctrlruntime.GetConfigOrDie()

		qps, burst := *rateLimiterQPS, *rateLimiterBurst
		if qps != 0 {
			cfg.QPS = float32(qps)
		}
		if burst != 0 {
			cfg.Burst = burst
		}
		if burst != 0 && qps != 0 {
			setupLog.Info("Configuring rate limiter", "QPS", qps, "burst", burst)
			cfg.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(cfg.QPS, cfg.Burst)
		}

		managerOpts.KubeConfig = cfg
	}

	if *profilerAddress != "" {
		setupLog.Info(
			"Profiler listening for requests",
			"profiler-address", *profilerAddress)
		go runProfiler(*profilerAddress)
	}

	setupLog.Info("wait for webhook certificates")
	waitForWebhookCertificates(setupLog, managerOpts)

	// Create a function that adds all of the controllers and webhooks to the manager.
	addToManager := func(ctx *context.ControllerManagerContext, mgr ctrlmgr.Manager) error {
		if err := controllers.AddToManager(ctx, mgr); err != nil {
			return err
		}

		if lib.IsVMServiceV1Alpha2FSSEnabled() {
			if err := addConversionWebhooksToManager(ctx, mgr); err != nil {
				return err
			}
		}

		return webhooks.AddToManager(ctx, mgr)
	}

	setupLog.Info("creating controller manager")
	managerOpts.InitializeProviders = manager.InitializeProviders
	managerOpts.AddToManager = addToManager
	mgr, err := manager.New(managerOpts)
	if err != nil {
		setupLog.Error(err, "problem creating controller manager")
		os.Exit(1)
	}

	setupLog.Info("setting up webhook server TLS config")
	webhookServer := mgr.GetWebhookServer()
	srv := webhookServer.(*webhook.DefaultServer)
	configureWebhookTLS(&srv.Options)

	setupLog.Info("adding readiness check to controller manager")
	if err := mgr.AddReadyzCheck("webhook", webhookServer.StartedChecker()); err != nil {
		setupLog.Error(err, "unable to create readiness check")
		os.Exit(1)
	}

	setupLog.Info("starting controller manager")
	sigHandler := ctrlsig.SetupSignalHandler()
	if err := mgr.Start(sigHandler); err != nil {
		setupLog.Error(err, "problem running controller manager")
		os.Exit(1)
	}
}

// addConversionWebhooksToManager adds the ctrl-runtime managed webhooks. We just use these
// for version conversion, but they can also do mutation and validation webhook callbacks
// instead of our separate webhooks.
func addConversionWebhooksToManager(_ *context.ControllerManagerContext, mgr ctrlmgr.Manager) error {
	if err := (&v1alpha1.VirtualMachine{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha1.VirtualMachineClass{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha1.VirtualMachineImage{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha1.ClusterVirtualMachineImage{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha1.VirtualMachinePublishRequest{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha1.VirtualMachineService{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha1.VirtualMachineSetResourcePolicy{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := (&v1alpha2.VirtualMachine{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha2.VirtualMachineClass{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha2.VirtualMachineImage{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha2.ClusterVirtualMachineImage{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha2.VirtualMachinePublishRequest{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha2.VirtualMachineService{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&v1alpha2.VirtualMachineSetResourcePolicy{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	return nil
}

func configureWebhookTLS(opts *webhook.Options) {
	tlsCfgFunc := func(cfg *tls.Config) {
		cfg.MinVersion = tls.VersionTLS12
		cfg.CipherSuites = []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		}
	}
	opts.TLSOpts = []func(*tls.Config){
		tlsCfgFunc,
	}
}

func waitForWebhookCertificates(setupLog logr.Logger, managerOpts manager.Options) {
	waitOnCertsStartTime := time.Now()
	for {
		select {
		case <-certDirReady(managerOpts.WebhookSecretVolumeMountPath):
			return
		case <-time.After(time.Second * 5):
			setupLog.Info("waiting on certificates", "elapsed", time.Since(waitOnCertsStartTime).String())
		}
	}
}

// CertDirReady returns a channel that is closed when there are certificates
// available in the configured certificate directory. If CertDir is
// empty or the specified directory does not exist, then the returned channel
// is never closed.
func certDirReady(certDir string) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		crtPath := path.Join(certDir, serverCertName)
		keyPath := path.Join(certDir, serverKeyName)
		for {
			if file, err := os.Stat(crtPath); err == nil {
				if file.Size() > 0 {
					if file, err := os.Stat(keyPath); err == nil {
						if file.Size() > 0 {
							close(done)
							return
						}
					}
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()
	return done
}

func runProfiler(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	_ = http.ListenAndServe(addr, mux)
}
