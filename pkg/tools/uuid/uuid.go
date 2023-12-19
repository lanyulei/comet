package uuid

import (
	"strings"

	"github.com/google/uuid"
)

/*
  @Author : lanyulei
  @Desc :
*/

func Get() (uid string) {
	uuidWithHyphen := uuid.New()
	uid = strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return
}
