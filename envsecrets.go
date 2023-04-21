package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"strings"
	"log"
	"os"
	"os/exec"
	"context"
    "google.golang.org/api/option"
    "google.golang.org/api/secretmanager/v1"
)
var fileConfig string = "./envsecrets-config.json"

//Get input file json with env and name of secrets to read
func getInputSecretsJSON() map[string]string {
	secrets := make(map[string]string)
	jsonFile, err := os.Open(fileConfig)
	if err != nil {
		log.Fatalf(fmt.Sprintf("%s",err))
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	if !gjson.Valid(string(byteValue)) {
		log.Fatalf(fmt.Sprintf("Invalid input file %s",fileConfig))
	}

	value := gjson.Get(string(byteValue), "secrets")
	if !value.Exists() {
		log.Fatalf(fmt.Sprintf("The input file %s don't contains secrets",fileConfig))	
	}	
	value.ForEach(func(key, value gjson.Result) bool {

		item := gjson.GetMany(value.String(), "env", "name")
		env := item[0].String()
		name := item[1].String()
		secrets[env] = name
		return true
	})
	return secrets	
}

//Get input file json with configuration
func getInputConfigJSON() (bool) {
    jsonString, err := ioutil.ReadFile(fileConfig)
    if err != nil {
		log.Fatalf(fmt.Sprintf("%s",err))
    }
	convert_to_uppercase_var_names := gjson.Get(string(jsonString), "config.convert_to_uppercase_var_names")
	if !convert_to_uppercase_var_names.Exists() {
		return false
	}	
	return convert_to_uppercase_var_names.Bool()
}

// Get secret payload data from Secret Manager
func getGCPSecretManager(name string) string {
	// Crete a backgound context
	ctx := context.Background()
 
	// Crear un cliente de Secret Manager utilizando las credenciales activas del SDK de GCP
	client, err := secretmanager.NewService(ctx, option.WithCredentialsFile(""), option.WithUserAgent("c-army/1.0"))
	if err != nil {
		fmt.Println("If you are running this in GKE Verify the Workload Identity setup https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#verify_the_setup")
		log.Fatalf(fmt.Sprintf("%s",err))
		//error could not find default credentials. See https://cloud.google.com/docs/authentication/external/set-up-adc for more information
	}
  
	// Get secret version
	result, err := client.Projects.Secrets.Versions.Access(name).Do()
	if err != nil {
		panic(err)
	}
 
	// Get secret value
	payload, err := base64.StdEncoding.DecodeString(result.Payload.Data)
	if err != nil {
		panic(err)
	}
   return string(payload)
 }


func main() {
    configPath := os.Getenv("ENVSECRETS_CONFIG_PATH")
	if configPath != "" {
		fileConfig = configPath
	}
	// check argumnets
	args := os.Args
	if len(args) == 1 {
		log.Fatalf("Error: An argument was expected, for example, ./envsecrets ./entrypoint.sh")
	}
	cmd := exec.Command(args[1], args[2:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	convert_to_uppercase_var_names := getInputConfigJSON()
	secrets := getInputSecretsJSON()
	for env, name := range secrets {
		secretPayload := getGCPSecretManager(name)
		//fmt.Println(env)
		//fmt.Println(secretPayload)

        if gjson.Valid(secretPayload) {
			//fmt.Println("el secret es un json")
			// Iterar sobre las claves y valores de la estructura JSON
			result := gjson.Parse(secretPayload)
			result.ForEach(func(key, value gjson.Result) bool {
				//fmt.Printf("key: %s, value: %s\n", key.String(), value.String())
				var_name := key.String()
				if convert_to_uppercase_var_names {
					var_name = strings.ToUpper(var_name)
				}
				var_name = strings.Replace(var_name, "-", "_", 1)
				cmd.Env = append(cmd.Env, var_name+"="+value.String())	
				return true // Continuar iterando
			})
		}else {
			cmd.Env = append(cmd.Env, env+"="+secretPayload)
		}
	}
	//create env secrets
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}	
}
