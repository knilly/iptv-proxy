/*
 * Iptv-Proxy is a project to proxyfie an m3u file and to proxyfie an Xtream iptv service (client API).
 * Copyright (C) 2020  Pierre-Emmanuel Jacquier
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/pierre-emmanuelJ/iptv-proxy/pkg/config"
	"github.com/pierre-emmanuelJ/iptv-proxy/pkg/regex"
	"github.com/pierre-emmanuelJ/iptv-proxy/pkg/server"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iptv-proxy",
	Short: "Reverse proxy on iptv m3u file and xtream codes server api",
	Run: func(cmd *cobra.Command, args []string) {
		confRegex := &regex.RegexSettings{
			M3uGroup:       viper.GetString("regex-m3u-group"),
			M3uChannel:     viper.GetString("regex-m3u-channel"),
			CategoryLive:   viper.GetString("regex-category-live"),
			CategoryVod:    viper.GetString("regex-category-vod"),
			CategorySeries: viper.GetString("regex-category-series"),
			StreamsLive:    viper.GetString("regex-streams-live"),
			StreamsSeries:  viper.GetString("regex-streams-series"),
			StreamsVod:     viper.GetString("regex-streams-vod"),
		}

		regex, regexErr := regex.NewFilter(confRegex)
		if regexErr != nil {
			log.Fatal(regexErr)
		}

		m3uURL := viper.GetString("m3u-url")
		remoteHostURL, err := url.Parse(m3uURL)
		if err != nil {
			log.Fatal(err)
		}

		xtreamUser := viper.GetString("xtream-user")
		xtreamPassword := viper.GetString("xtream-password")
		xtreamBaseURL := viper.GetString("xtream-base-url")

		var username, password string
		if strings.Contains(m3uURL, "/get.php") {
			username = remoteHostURL.Query().Get("username")
			password = remoteHostURL.Query().Get("password")
		}

		if xtreamBaseURL == "" && xtreamPassword == "" && xtreamUser == "" {
			if username != "" && password != "" {
				log.Printf("[iptv-proxy] INFO: It's seams you are using an Xtream provider!\n")

				xtreamUser = username
				xtreamPassword = password
				xtreamBaseURL = fmt.Sprintf("%s://%s", remoteHostURL.Scheme, remoteHostURL.Host)
				log.Printf("[iptv-proxy] INFO: xtream service enable with xtream base url: %q xtream username: %q xtream password: %q\n", xtreamBaseURL, xtreamUser, xtreamPassword)
			}
		}

		conf := &config.ProxyConfig{
			HostConfig: &config.HostConfiguration{
				Hostname: viper.GetString("hostname"),
				Port:     viper.GetInt("port"),
			},
			RemoteURL:            remoteHostURL,
			XtreamUser:           config.CredentialString(xtreamUser),
			XtreamPassword:       config.CredentialString(xtreamPassword),
			XtreamBaseURL:        xtreamBaseURL,
			M3UCacheExpiration:   viper.GetInt("m3u-cache-expiration"),
			User:                 config.CredentialString(viper.GetString("user")),
			Password:             config.CredentialString(viper.GetString("password")),
			AdvertisedPort:       viper.GetInt("advertised-port"),
			HTTPS:                viper.GetBool("https"),
			M3UFileName:          viper.GetString("m3u-file-name"),
			CustomEndpoint:       viper.GetString("custom-endpoint"),
			CustomId:             viper.GetString("custom-id"),
			XtreamGenerateApiGet: viper.GetBool("xtream-api-get"),

			Filter: regex,
		}

		if conf.AdvertisedPort == 0 {
			conf.AdvertisedPort = conf.HostConfig.Port
		}

		server, err := server.NewServer(conf)
		if err != nil {
			log.Fatal(err)
		}

		if e := server.Serve(); e != nil {
			log.Fatal(e)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "iptv-proxy-config", "C", "Config file (default is $HOME/.iptv-proxy.yaml)")
	rootCmd.Flags().StringP("m3u-url", "u", "", `Iptv m3u file or url e.g: "http://example.com/iptv.m3u"`)
	rootCmd.Flags().StringP("m3u-file-name", "", "iptv.m3u", `Name of the new proxified m3u file e.g "http://poxy.com/iptv.m3u"`)
	rootCmd.Flags().StringP("custom-endpoint", "", "", `Custom endpoint "http://poxy.com/<custom-endpoint>/iptv.m3u"`)
	rootCmd.Flags().StringP("custom-id", "", "", `Custom anti-collison ID for each track "http://proxy.com/<custom-id>/..."`)
	rootCmd.Flags().Int("port", 8080, "Iptv-proxy listening port")
	rootCmd.Flags().Int("advertised-port", 0, "Port to expose the IPTV file and xtream (by default, it's taking value from port) useful to put behind a reverse proxy")
	rootCmd.Flags().String("hostname", "", "Hostname or IP to expose the IPTVs endpoints")
	rootCmd.Flags().BoolP("https", "", false, "Activate https for urls proxy")
	rootCmd.Flags().String("user", "usertest", "User auth to access proxy (m3u/xtream)")
	rootCmd.Flags().String("password", "passwordtest", "Password auth to access proxy (m3u/xtream)")
	rootCmd.Flags().String("xtream-user", "", "Xtream-code user login")
	rootCmd.Flags().String("xtream-password", "", "Xtream-code password login")
	rootCmd.Flags().String("xtream-base-url", "", "Xtream-code base url e.g(http://expample.tv:8080)")
	rootCmd.Flags().Int("m3u-cache-expiration", 1, "M3U cache expiration in hour")
	rootCmd.Flags().BoolP("xtream-api-get", "", false, "Generate get.php from xtream API instead of get.php original endpoint")

	// Regex options
	rootCmd.Flags().String("regex-m3u-group", "", "Regex applied to filtering M3U groups")
	rootCmd.Flags().String("regex-m3u-channel", "", "Regex applied to filtering M3U channels")
	rootCmd.Flags().String("regex-category-live", "", "Regex applied to filtering live categories")
	rootCmd.Flags().String("regex-category-vod", "", "Regex applied to filtering vod categories")
	rootCmd.Flags().String("regex-category-series", "", "Regex applied to filtering series categories")
	rootCmd.Flags().String("regex-streams-live", "", "Regex applied to filtering live streams")
	rootCmd.Flags().String("regex-streams-vod", "", "Regex applied to filtering vod streams")
	rootCmd.Flags().String("regex-streams-series", "", "Regex applied to filtering series streams")

	if e := viper.BindPFlags(rootCmd.Flags()); e != nil {
		log.Fatal("error binding PFlags to viper")
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".iptv-proxy" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".iptv-proxy")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
