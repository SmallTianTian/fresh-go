package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/test"
	"github.com/SmallTianTian/fresh-go/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewProject(t *testing.T) {
	Convey("", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir)

		config.DefaultConfig.Project.Path = dir
		config.DefaultConfig.Project.Name = "fresh"
		config.DefaultConfig.Project.Remote = "github.com"
		config.DefaultConfig.Project.Vendor = false

		Convey("能正常初始化项目", func() {
			proPath := filepath.Join(dir, config.DefaultConfig.Project.Name)

			So(utils.IsExist(proPath), ShouldBeFalse)

			So(NewProject, ShouldNotPanic)
			So(utils.IsExist(proPath), ShouldBeTrue)

			mod := filepath.Join(proPath, "go.mod")
			vendor := filepath.Join(proPath, "vendor")
			So(utils.IsExist(mod), ShouldBeTrue)
			So(utils.IsExist(vendor), ShouldBeFalse)

			remote, _, name := utils.GetRemoteOwnerAndProjectName(proPath)
			So(remote, ShouldEqual, config.DefaultConfig.Project.Remote)
			So(name, ShouldEqual, config.DefaultConfig.Project.Name)

			dir := test.TempDir()
			defer os.RemoveAll(dir)

			config.DefaultConfig.Project.Path = dir
			So(func() { utils.Exec(proPath, "go", "vet") }, ShouldNotPanic)
			So(utils.Exec(proPath, "go", "vet", "./..."), ShouldBeNil)
		})

		Convey("能正常初始化包含 vendor 的项目", func() {
			config.DefaultConfig.Project.Vendor = true

			proPath := filepath.Join(dir, config.DefaultConfig.Project.Name)

			So(utils.IsExist(proPath), ShouldBeFalse)

			So(NewProject, ShouldNotPanic)
			So(utils.IsExist(proPath), ShouldBeTrue)

			mod := filepath.Join(proPath, "go.mod")
			vendor := filepath.Join(proPath, "vendor")
			So(utils.IsExist(mod), ShouldBeTrue)
			So(utils.IsExist(vendor), ShouldBeTrue)

			remote, _, name := utils.GetRemoteOwnerAndProjectName(proPath)
			So(remote, ShouldEqual, config.DefaultConfig.Project.Remote)
			So(name, ShouldEqual, config.DefaultConfig.Project.Name)

			dir := test.TempDir()
			defer os.RemoveAll(dir)

			config.DefaultConfig.Project.Path = dir
			So(func() { utils.Exec(proPath, "go", "vet") }, ShouldNotPanic)
			So(utils.Exec(proPath, "go", "vet", "./..."), ShouldBeNil)
		})

		Convey("不可重复初始化项目", func() {
			config.DefaultConfig.Project.Vendor = true
			proPath := filepath.Join(dir, config.DefaultConfig.Project.Name)

			test.InitGoMod("old", "", "test", proPath)
			mod := filepath.Join(proPath, "go.mod")

			So(utils.IsExist(proPath), ShouldBeTrue)
			So(utils.IsExist(mod), ShouldBeTrue)

			So(NewProject, ShouldPanic)
		})
	})
}
