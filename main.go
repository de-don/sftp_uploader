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
	"github.com/kardianos/osext"
)

func generateImageName(format string) string {
	// generate screen name
	t := time.Now()
	return fmt.Sprintf("%s.png", t.Format(format))
}

func loadConfig(cfgName string) (map[string]string, error) {
	executablePath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal("Error: Couldn't determine working directory: " + err.Error())
	}
	// Set the working directory to the path the executable is located in.
	os.Chdir(executablePath)


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

func notify(title, text string) {
	log.Println(text)
	if err := beeep.Notify(title, text, ""); err != nil{
		panic(err)
	}
}

func alertError(title, text string) {
	log.Println(text)
	if err := beeep.Alert(title, text, ""); err != nil{
		panic(err)
	}
	os.Exit(1)
}

func main() {
	var res = flag.Bool("r", false, "true - load from clipboard, 1 - load from stdin")
	flag.Parse()


	// load configuration
	config, err := loadConfig("./config.ini")
	if err != nil {
		alertError("Config error", "Fail to read file: " + err.Error())
	}

	// generate screen name
	screenName := generateImageName(config["nameFormat"])

	// get raw image from clipboard or stdin
	var image []byte
	if *res == false {
		if image, err = getImageFromClipboard(); err != nil {
			alertError("Image error", "Image not found in buffer")
		}
	} else {
		if image, err = getImageFromStdin(); err != nil {
			alertError("Image error", "Image not found in stdin")
		}
	}

	// check that it is really image
	if !isImage(image){
		alertError("Image error", "Input data not is image")
	}

	// create ssh connection
	connection, err := createSshSession(
		config["host"],
		config["port"],
		config["username"],
		config["password"],
	)
	if err != nil {
		alertError("Connection error", "Failed to dial: " + err.Error())
	}

	// create sftp connection
	sftpConn, err := sftp.NewClient(connection)
	if err != nil {
		alertError("Connection error", err.Error())
	}
	defer sftpConn.Close()

	// upload image
	if saveImageOnServer(sftpConn, screenName, image) != nil {
		log.Fatalf("Can't upload image: %s", err)
		alertError("Connection error", "Can't upload image: " + err.Error())
	}

	// save link to clipboard
	link := fmt.Sprintf("http://%s/%s/%s", config["host"], config["username"], screenName)
	putTextToClipboard(link)

	notify("Image uploaded", link)
}
