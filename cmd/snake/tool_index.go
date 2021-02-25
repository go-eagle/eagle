package main

import "time"

var toolIndexs = []*Tool{
	{
		Name:      "snake",
		Alias:     "snake",
		BuildTime: time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/1024casts/snake/cmd/snake@" + Version,
		Summary:   "Snake工具集",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "snake",
		Hidden:    true,
	},
	{
		Name:      "protoc",
		Alias:     "gen-protoc",
		BuildTime: time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/1024casts/snake/cmd/gen-protoc@" + Version,
		Summary:   "快速方便生成pb.go的protoc封装，windows、Linux请先安装protoc工具",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "snake",
	},
	{
		Name:         "gen-project",
		Alias:        "snake-gen-project",
		Install:      "go get -u github.com/1024casts/snake/cmd/snake-gen-project",
		BuildTime:    time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local),
		Platform:     []string{"darwin", "linux", "windows"},
		Hidden:       true,
		Requirements: []string{"wire"},
	},
	//  third party
	{
		Name:      "wire",
		Alias:     "wire",
		BuildTime: time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/google/wire/cmd/wire",
		Platform:  []string{"darwin", "linux", "windows"},
		Hidden:    true,
	},
	{
		Name:      "swagger",
		Alias:     "swagger",
		BuildTime: time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/go-swagger/go-swagger/cmd/swagger",
		Summary:   "swagger api文档",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "goswagger.io",
	},
}
