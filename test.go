package king

import (
	sh "github.com/codeskyblue/go-sh"
	"fmt"
)

func main(){
	output, err := sh.Command("sh", "shells/mvn.sh").Output()
	fmt.Println(string(output), err)
}
