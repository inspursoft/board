package service

import (
	"git/inspursoft/board/src/adminserver/models"
	"regexp"
	"strings"
)

//GetMonitor returns Board containers' information.
func GetMonitor() ([]*models.Boardinfo, error) {
	command := `docker ps -a --format "table {{.ID}}\t{{.Image}}\t{{.CreatedAt}}\t{{.Status}}\t{{.Ports}}" | grep -v adminserver | grep ` + models.ImagePrefix
	resp, err := Execute(command)
	if err != nil {
		return nil, err
	}
	row := strings.Count(resp, "\n")
	arr := strings.Split(resp, "\n")

	//var containers [row]*models.Boardinfo
	containersAdd := make([]*models.Boardinfo, row)
	containersVal := make([]models.Boardinfo, row)

	containerPre := models.ContainerPrefix
	_, devopsOpt, err := GetFileFromDevopsOpt()
	if err != nil {
		return nil, err
	}
	if devopsOpt == "legacy" {
		containerPre = "archive"
	}
	command2 := `docker stats -a --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}\t{{.MemPerc}}\t{{.PIDs}}" | grep -v adminserver | grep ` + containerPre
	resp2, err := Execute(command2)
	if err != nil {
		return nil, err
	}
	arr2 := strings.Split(resp2, "\n")

	reg, _ := regexp.Compile("\\s{2,}")
	for i := 0; i < row; i++ {
		items := reg.Split(arr[i], -1)

		//assign port with null if missing.
		if len(items) < 5 {
			items = append(items, "")
		}

		itemsStats := reg.Split(arr2[i], -1)
		items = append(items, itemsStats...)

		containersVal[i] = models.GetBoardinfo(items)
		containersAdd[i] = &containersVal[i]
		//fmt.Printf("%q\n", items)
	}

	return containersAdd, nil
}
