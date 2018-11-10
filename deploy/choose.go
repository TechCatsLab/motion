package deploy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var commandsList map[string][]string

func init() {
	commandsList = make(map[string][]string)

	file, err := os.Open("./commands")
	if err != nil {
		panic(err)
	}
	li, err := file.Readdir(-1)
	if err != nil {
		panic(err)
	}

	for _, each := range li {
		if !each.IsDir() {
			continue
		}

		eaea, err := ioutil.ReadDir("./commands/" + each.Name())
		if err != nil {
			panic(err)
		}

		commandsList[each.Name()] = make([]string, 0, len(eaea))
		for _, e := range eaea {
			if e.IsDir() {
				continue
			}
			commandsList[each.Name()] = append(commandsList[each.Name()], e.Name())
		}
	}
}

func ChooseOS() {
	var curOS int = -1

	archs := make([]string, 0, len(commandsList))
	for key, _ := range commandsList {
		archs = append(archs, key)
	}
	for i, n := range archs {
		fmt.Printf("%d) %s\n", i, n)
	}
	fmt.Println("Choose os: ")
	fmt.Scanf("%d", &curOS)
	if curOS < 0 || curOS >= len(archs) {
		return
	}

	EnvConfig.OS = archs[curOS]
}

func ChooseOperations() []string {
	var choice string

	fmt.Printf("Your OS is %s\n\nChoose operations\n", EnvConfig.OS)
	ops := commandsList[EnvConfig.OS]
	for i, n := range ops {
		fmt.Printf("%c) %s\n", i+97, n[:len(n)-5])
	}
	fmt.Scanln(&choice)

	choice = strings.ToLower(choice)
	result := make([]string, 0, 3)
	for _, ch := range choice {
		i := int(ch) - 97
		if i < 0 || i > len(ops)-1 {
			continue
		}
		result = append(result, ops[i])
	}

	return result
}

func readCommands(o, filename string) ([]string, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("./commands/%s/%s", o, filename))
	if err != nil {
		return nil, err
	}
	var x interface{}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return nil, err
	}
	list := x.(map[string]interface{})["commands"].([]interface{})
	res := make([]string, 0, len(list))
	for _, c := range list {
		res = append(res, c.(string))
	}
	return res, nil
}

func RunCommands(file string) error {
	commands, err := readCommands(EnvConfig.OS, file)
	if err != nil {
		return err
	}
	length := len(commands)
	for i := 0; i < length; i++ {
		err = run(commands[i])
		if err != nil {
			return err
		}
	}
	return nil
}
