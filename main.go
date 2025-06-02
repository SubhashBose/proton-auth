package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	proton "github.com/mort666/go-proton-api"
	common "rtlabs.cloud/protonsession"
)

func main() {
	rawOutput := flag.Bool("raw", false, "Print only UID and AccessToken, no labels")
	hostURL := flag.String("hosturl", "https://mail.proton.me/api", "Override the default Proton API host URL")

	flag.Parse()

	ctx := context.Background()

	username := os.Getenv("PROTON_USERNAME")
	password := os.Getenv("PROTON_PASSWORD")

	if username == "" || password == "" {
		fmt.Fprintln(os.Stderr, "Error: PROTON_USERNAME and PROTON_PASSWORD must be set as environment variables to get auth tokens.")
		os.Exit(1)
	}

	protonOptions := []proton.Option{
		proton.WithAppVersion("other"),
		proton.WithHostURL(*hostURL),
	}

	sessionStore := common.NewFileStore("proton-sessions.db", "default")
	sessionStore.CacheDir = false

	var pmSession *common.Session
	sessionConfig, err := sessionStore.Load()
	if err == nil {
		sessionCreds := &common.SessionCredentials{
			UID:          sessionConfig.UID,
			AccessToken:  sessionConfig.AccessToken,
			RefreshToken: sessionConfig.RefreshToken,
		}
		pmSession, err = common.SessionFromRefresh(ctx, protonOptions, sessionCreds)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SessionFromRefresh failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		pmSession, err = common.SessionFromLogin(ctx, protonOptions, username, password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SessionFromLogin failed: %v\n", err)
			os.Exit(1)
		}
	}

	if *rawOutput {
		fmt.Println(pmSession.Auth.AccessToken)
		fmt.Println(pmSession.Auth.UID)
	} else {
		fmt.Println("AccessToken:", pmSession.Auth.AccessToken)
		fmt.Println("UID:", pmSession.Auth.UID)
	}
}
