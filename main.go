package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/pkg/sftp"
	"log"
	"os"
	"time"
	"flag"
	"github.com/gen2brain/beeep"
	"runtime"
	"path"
)

func generateImageName(format string) string {
	// generate screen name
	t := time.Now()
	return fmt.Sprintf("%s.png", t.Format(format))
}

func loadConfig(cfgName string) (map[string]string, error) {
	_, filename, _, _ := runtime.Caller(1)
	cfgName = path.Join(path.Dir(filename), cfgName)


	// load configuration
	cfg, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment:         true,
		UnescapeValueCommentSymbols: true,
	}, cfgName)

	if err != nil {
		return nil, err
	}
	data := map[string]string{
		"username":   cfg.Section("credentials").Key("username").String(),
		"password":   cfg.Section("credentials").Key("password").String(),
		"host":       cfg.Section("server").Key("host").String(),
		"port":       cfg.Section("server").Key("port").String(),
		"nameFormat": cfg.Section("other").Key("name_pattern").String(),
	}
	return data, nil
}

func main() {
	var res = flag.Bool("r", false, "true - load from clipboard, 1 - load from stdin")
	flag.Parse()


	// load configuration
	config, err := loadConfig("./config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// generate screen name
	screenName := generateImageName(config["nameFormat"])

	var image []byte
	if *res == false {
		// get image
		image, err = getImageFromClipboard()
		if err != nil {
			log.Fatal("Image not found in buffer")
			os.Exit(1)
		}
	} else {
		image, err = getImageFromStdin()
	}

	if !isImage(image){
		println(image)
		log.Fatal("Not a image")
		os.Exit(1)
	}

	// create ssh connection
	connection, err := createSshSession(
		config["host"],
		config["port"],
		config["username"],
		config["password"],
	)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}

	// create sftp connection
	sftpConn, err := sftp.NewClient(connection)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpConn.Close()

	if saveImageOnServer(sftpConn, screenName, image) != nil {
		log.Fatalf("Can't upload image: %s", err)
	}

	// save link to clipboard
	link := fmt.Sprintf("http://%s/%s/%s", config["host"], config["username"], screenName)
	log.Println(link)
	putTextToClipboard(link)
	err = beeep.Notify("Uploaded", link, "assets/information.png")
	if err != nil {
		panic(err)
	}
}
