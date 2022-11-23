package main

import "github.com/spf13/viper"

// loadEnvConfiguration loads environment variables
func loadEnvConfiguration() {
	// HTTP server
	viper.SetDefault("HTTPBindAddress", "localhost:8080")
	_ = viper.BindEnv("HTTPBindAddress", "HTTP_BIND_ADDRESS")

	viper.SetDefault("ExternalLocationAddress", "http://localhost:8080")
	_ = viper.BindEnv("ExternalLocationAddress", "EXTERNAL_LOCATION_ADDRESS")

	viper.SetDefault("HTTPShutdownTimeoutInSeconds", "30")
	_ = viper.BindEnv("HTTPShutdownTimeoutInSeconds", "HTTP_SHUTDOWN_TIMEOUT_SECONDS")

	// External user service
	viper.SetDefault("UserServiceGRPCDialTarget", "localhost:50051")
	_ = viper.BindEnv("UserServiceGRPCDialTarget", "USER_SERVICE_GRPC_DIAL_TARGET")

}
