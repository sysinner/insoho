// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"fmt"
	"sort"

	"github.com/hooto/iam/iamapi"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/skv"

	incfg "github.com/sysinner/incore/config"
	"github.com/sysinner/incore/inapi"
)

var (
	init_sys_user    = "sysadmin"
	SysConfigurators = []*inapi.SysConfigurator{
		{
			Name:  "innerstack/sys/webui",
			Title: "WebUI",
			Fields: []*inapi.AppConfigField{
				{
					Name:  "html_head_title",
					Title: "Web Html Head Title",
					Type:  inapi.AppConfigFieldTypeString,
				},
				{
					Name:  "cp_navbar_logo",
					Title: "Logo of Panel",
					Type:  inapi.AppConfigFieldTypeString,
				},
				{
					Name:  "cp_navbar_title",
					Title: "Navbar Title of Panel",
					Type:  inapi.AppConfigFieldTypeString,
				},
				{
					Name:  "ops_navbar_title",
					Title: "Navbar Title of Ops Panel",
					Type:  inapi.AppConfigFieldTypeString,
				},
				/**
				{
					Name:  "page_footer_copyright",
					Title: "Page Footer Copyright Information",
					Type:  inapi.AppConfigFieldTypeString,
				},
				*/
			},
			ReadRoles: []uint32{100},
		},
	}
)

func InitHostletData() map[string]interface{} {

	items := map[string]interface{}{}

	return items
}

func InitIamAccessKeyData() []iamapi.AccessKey {
	return []iamapi.AccessKey{
		init_cache_akacc,
	}
}

var (
	init_zmd_items         = map[string]interface{}{}
	init_zmd_items_upgrade = map[string]interface{}{}
)

func InitZoneMasterData() map[string]interface{} {

	var (
		name          string
		host_id       = incfg.Config.Host.Id
		host_lan_addr = string(incfg.Config.Host.LanAddr)
		host_wan_addr = string(incfg.Config.Host.WanAddr)
	)

	// sys zone
	sys_zone := inapi.ResZone{
		Meta: &inapi.ObjectMeta{
			Id:      incfg.Config.Host.ZoneId,
			Name:    incfg.Config.Host.ZoneId,
			Created: uint64(types.MetaTimeNow()),
			Updated: uint64(types.MetaTimeNow()),
		},
		Phase: 1,
		LanAddrs: []string{
			host_lan_addr,
		},
	}
	if incfg.Config.Host.WanAddr.Valid() {
		sys_zone.WanAddrs = []string{
			host_wan_addr,
		}
	}

	//
	sys_zone.OptionSet("iam/acc_charge/access_key", init_cache_akacc.AccessKey)
	sys_zone.OptionSet("iam/acc_charge/secret_key", init_cache_akacc.SecretKey)

	init_zmd_items[inapi.NsGlobalSysZone(incfg.Config.Host.ZoneId)] = sys_zone
	init_zmd_items[inapi.NsZoneSysInfo(incfg.Config.Host.ZoneId)] = sys_zone

	// inapi.ObjPrint("sys_zone", sys_zone)

	// sys cell
	sys_cell := inapi.ResCell{
		Meta: &inapi.ObjectMeta{
			Id:      incfg.Config.Host.CellId,
			Name:    incfg.Config.Host.CellId,
			Created: uint64(types.MetaTimeNow()),
			Updated: uint64(types.MetaTimeNow()),
		},
		ZoneId:      incfg.Config.Host.ZoneId,
		Phase:       1,
		Description: "",
	}
	init_zmd_items[inapi.NsGlobalSysCell(incfg.Config.Host.ZoneId, incfg.Config.Host.CellId)] = sys_cell
	init_zmd_items[inapi.NsZoneSysCell(incfg.Config.Host.ZoneId, incfg.Config.Host.CellId)] = sys_cell

	// sys host
	sys_host := inapi.ResHost{
		Meta: &inapi.ObjectMeta{
			Id:      host_id,
			Created: uint64(types.MetaTimeNow()),
			Updated: uint64(types.MetaTimeNow()),
		},
		Operate: &inapi.ResHostOperate{
			Action: 1,
			ZoneId: incfg.Config.Host.ZoneId,
			CellId: incfg.Config.Host.CellId,
		},
		Spec: &inapi.ResHostSpec{
			PeerLanAddr: host_lan_addr,
		},
	}
	init_zmd_items[inapi.NsZoneSysHost(incfg.Config.Host.ZoneId, host_id)] = sys_host
	init_zmd_items[inapi.NsZoneSysHostSecretKey(incfg.Config.Host.ZoneId, host_id)] = incfg.Config.Host.SecretKey

	//

	// zone-master node(s)/leader
	init_zmd_items[inapi.NsZoneSysMasterNode(incfg.Config.Host.ZoneId, host_id)] = inapi.ResZoneMasterNode{
		Id:     host_id,
		Addr:   host_lan_addr,
		Action: 1,
	}
	// init_zmd_items[inapi.NsKvZoneSysMasterLeader(incfg.Config.Host.ZoneId)] = host_id

	//
	name = "g1"
	plan_g1 := inapi.PodSpecPlan{
		Meta: types.InnerObjectMeta{
			ID:      name,
			Name:    "General g1",
			User:    "sysadmin",
			Version: "1",
			Created: types.MetaTimeNow(),
			Updated: types.MetaTimeNow(),
		},
		Status: inapi.SpecStatusActive,
		Zones: []*inapi.PodSpecPlanZoneBound{
			{
				Name:  incfg.Config.Host.ZoneId,
				Cells: types.ArrayString([]string{incfg.Config.Host.CellId}),
			},
		},
	}
	plan_g1.Labels.Set("pod/spec/plan/type", "g1")
	plan_g1.Annotations.Set("meta/homepage", "http://www.sysinner.com")
	plan_g1.SortOrder = 2

	//
	plan_t1 := plan_g1
	plan_t1.Meta.ID = "t1"
	plan_t1.Meta.Name = "Tiny t1"
	plan_t1.Labels.Set("pod/spec/plan/type", "t1")
	plan_t1.Annotations.Set("meta/homepage", "http://www.sysinner.com")
	plan_t1.SortOrder = 1

	// Spec/Image
	plan_g1.ImageDefault = inapi.BoxImageRepoDefault
	for _, vi := range [][]string{
		// {name, tag, driver, display-name}
		// {inapi.BoxImageRepoDefault, "a1el7v1", inapi.PodSpecBoxImageDocker, "General v1"},
		{inapi.BoxImageRepoDefault + "/innerstack-g2", "el7", inapi.PodSpecBoxImageDocker, "General v2"},
		// {inapi.BoxImageRepoDefault, "a2p1el7", inapi.PodSpecBoxImagePouch},
	} {

		image := inapi.PodSpecBoxImage{
			Meta: types.InnerObjectMeta{
				ID:      vi[0] + ":" + vi[1],
				Name:    vi[3],
				User:    "sysadmin",
				Version: "1",
				Created: types.MetaTimeNow(),
				Updated: types.MetaTimeNow(),
			},
			Name:      vi[0],
			Tag:       vi[1],
			Action:    inapi.PodSpecBoxImageActionEnable,
			Driver:    vi[2],
			SortOrder: 1,
			OsDist:    "el7",
			Arch:      inapi.SpecCpuArchAmd64,
		}
		init_zmd_items[string(inapi.NsGlobalBoxImage(vi[0], vi[1]))] = image

		image.Meta.User = ""
		image.Meta.Created = 0
		image.Meta.Updated = 0

		plan_g1.Images = append(plan_g1.Images, &inapi.PodSpecPlanBoxImageBound{
			RefId:    image.Meta.ID,
			RefName:  image.Name,
			RefTitle: image.Meta.Name,
			RefTag:   image.Tag,
			Driver:   image.Driver,
			OsDist:   image.OsDist,
			Arch:     image.Arch,
		})
	}

	plan_t1.Images = plan_g1.Images
	plan_t1.ImageDefault = plan_g1.ImageDefault

	for _, v := range [][]int32{
		//
		{1, 64}, // CPU 0.1 cores, RAM MB
		{1, 128},
		//
		{2, 128},
		{2, 256},
		//
		{4, 256},
		{4, 512},
		//
		{6, 512},
		{6, 1024},
		//
		{10, 512},
		{10, 1024},
		{10, 2048},
		{10, 4096},
		//
		{20, 1024},
		{20, 2048},
		{20, 4096},
		{20, 8192},
		//
		{40, 2048},
		{40, 4096},
		{40, 8192},
		{40, 16384},
		//
		{80, 4096},
		{80, 8192},
		{80, 16384},
		{80, 32768},
	} {

		name = fmt.Sprintf("c%dm%d", v[0], v[1])

		res := inapi.PodSpecResCompute{
			Meta: types.InnerObjectMeta{
				ID:      name,
				Name:    name,
				User:    "sysadmin",
				Created: types.MetaTimeNow(),
				Updated: types.MetaTimeNow(),
			},
			Status:   inapi.SpecStatusActive,
			CpuLimit: v[0],
			MemLimit: v[1],
		}

		if v[0] == 1 && v[1] == 128 {
			plan_t1.ResComputeDefault = res.Meta.ID
		}

		if v[0] == 10 && v[1] == 1024 {
			plan_g1.ResComputeDefault = res.Meta.ID
		}

		init_zmd_items[inapi.NsGlobalPodSpec("res/compute", res.Meta.ID)] = res

		res.Meta.User = ""
		res.Meta.Created = 0
		res.Meta.Updated = 0

		if v[0] < 10 {
			plan_t1.ResComputes = append(plan_t1.ResComputes, &inapi.PodSpecPlanResComputeBound{
				RefId:    res.Meta.ID,
				CpuLimit: res.CpuLimit,
				MemLimit: res.MemLimit,
			})
		}

		if v[0] >= 10 {
			plan_g1.ResComputes = append(plan_g1.ResComputes, &inapi.PodSpecPlanResComputeBound{
				RefId:    res.Meta.ID,
				CpuLimit: res.CpuLimit,
				MemLimit: res.MemLimit,
			})
		}
	}
	sort.Sort(plan_g1.ResComputes)
	sort.Sort(plan_t1.ResComputes)

	//
	vol_t1 := inapi.PodSpecResVolume{
		Meta: types.InnerObjectMeta{
			ID:      "lt1",
			Name:    "Tiny t1",
			User:    "sysadmin",
			Version: "1",
			Created: types.MetaTimeNow(),
			Updated: types.MetaTimeNow(),
		},
		Status:  inapi.SpecStatusActive,
		Limit:   10,
		Request: 1,
		Step:    1,
		Default: 1,
	}
	vol_t1.Labels.Set("pod/spec/res/volume/type", "system")

	vol_g1 := vol_t1
	vol_g1.Meta.ID = "lg1"
	vol_g1.Meta.Name = "General g1"
	vol_g1.Limit = 1000
	vol_g1.Request = 10
	vol_g1.Step = 10
	vol_g1.Default = 10

	vol_s1 := vol_t1
	vol_s1.Meta.ID = "ls1"
	vol_s1.Meta.Name = "SSD s1"
	vol_s1.Limit = 200
	vol_s1.Request = 10
	vol_s1.Step = 1
	vol_s1.Default = 10
	vol_s1.Attrs = inapi.ResVolValueAttrTypeSSD

	init_zmd_items[inapi.NsGlobalPodSpec("res/volume", vol_t1.Meta.ID)] = vol_t1
	init_zmd_items[inapi.NsGlobalPodSpec("res/volume", vol_g1.Meta.ID)] = vol_g1
	init_zmd_items[inapi.NsGlobalPodSpec("res/volume", vol_s1.Meta.ID)] = vol_s1

	vol_t1.Meta.User = ""
	vol_t1.Meta.Created = 0
	vol_t1.Meta.Updated = 0
	vol_g1.Meta.User = ""
	vol_g1.Meta.Created = 0
	vol_g1.Meta.Updated = 0

	//
	plan_t1.ResVolumes = append(plan_t1.ResVolumes, &inapi.PodSpecPlanResVolumeBound{
		RefId:   vol_t1.Meta.ID,
		RefName: vol_t1.Meta.Name,
		Limit:   vol_t1.Limit,
		Request: vol_t1.Request,
		Step:    vol_t1.Step,
		Default: vol_t1.Default,
	})
	plan_t1.ResVolumeDefault = vol_t1.Meta.ID

	//
	plan_g1.ResVolumes = append(plan_g1.ResVolumes, &inapi.PodSpecPlanResVolumeBound{
		RefId:   vol_g1.Meta.ID,
		RefName: vol_g1.Meta.Name,
		Limit:   vol_g1.Limit,
		Request: vol_g1.Request,
		Step:    vol_g1.Step,
		Default: vol_g1.Default,
	})
	plan_g1.ResVolumes = append(plan_g1.ResVolumes, &inapi.PodSpecPlanResVolumeBound{
		RefId:       vol_s1.Meta.ID,
		RefName:     vol_s1.Meta.Name,
		Limit:       vol_s1.Limit,
		Request:     vol_s1.Request,
		Step:        vol_s1.Step,
		Default:     vol_s1.Default,
		Attrs:       vol_s1.Attrs,
		ChargeValue: 0.0015,
	})

	plan_g1.ResVolumeDefault = vol_g1.Meta.ID

	//
	init_zmd_items[inapi.NsGlobalPodSpec("plan", plan_t1.Meta.ID)] = plan_t1
	init_zmd_items[inapi.NsGlobalPodSpec("plan", plan_g1.Meta.ID)] = plan_g1

	// specs := []string{
	// 	"app_spec_hooto-press-x1.json",
	// 	"app_spec_sysinner-httplb.json",
	// 	"app_spec_sysinner-mysql-x1.json",
	// }
	// for _, v := range specs {
	// 	var spec inapi.AppSpec
	// 	if err := json.DecodeFile(incfg.Prefix+"/misc/app-spec/"+v, &spec); err != nil || spec.Meta.ID == "" {
	// 		hlog.Printf("warn", "init app spec %s error", v)
	// 		continue
	// 	}
	// 	spec.Meta.User = "sysadmin"
	// 	init_zmd_items[inapi.NsGlobalAppSpec(spec.Meta.ID)] = spec
	// }

	for _, v := range SysConfigurators {
		incfg.SysConfigurators = append(incfg.SysConfigurators, v)
	}

	return init_zmd_items
}

func UpgradeZoneMasterData(data skv.Connector) error {

	if data == nil {
		return errors.New("UpgradeZoneMasterData skv.Connector Not Found")
	}

	return nil
}

func UpgradeIamData(data skv.Connector) error {

	if data == nil {
		return errors.New("UpgradeIamData skv.Connector Not Found")
	}

	return nil
}
