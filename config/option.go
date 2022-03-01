package config

import (
	"encoding/json"
	"github.com/scjtqs/jd_cookie/util"
	"os"
)

type Conf struct {
	DbConf       DbConf `json:"db_conf"`
	UpSave       string `json:"up_save"`
	UpSaveMethod string `json:"up_save_method"`
	UpSaveKey    string `json:"up_save_key"`
}

type DbConf struct {
	DbEnable   bool   `json:"db_enable"`
	DbType     string `json:"db_type"`
	DbHost     string `json:"db_host"`
	DbPort     string `json:"db_port"`
	DbUser     string `json:"db_user"`
	DbPass     string `json:"db_pass"`
	DbDatabase string `json:"db_database"`
}

// 通过路径获取配置信息
func GetConfigFronPath(c string) *Conf {
	conf := &Conf{}
	if !util.PathExists(c) {
		conf = defaultConf()
	} else {
		err := json.Unmarshal([]byte(util.ReadAllText(c)), conf)
		if err != nil {
			conf = defaultConf()
		}
	}
	return parseConfFromEnv(conf)
}

// 没有配置文件时候的默认配置
func defaultConf() *Conf {
	return &Conf{
		DbConf: DbConf{
			DbEnable: false,
			DbType:   "mysql",
		},
		UpSave: "",
	}
}

// 从环境变量中替换配置文件
func parseConfFromEnv(c *Conf) *Conf {
	if os.Getenv("UPSAVE") != "" {
		c.UpSave = os.Getenv("UPSAVE")
	}
	if os.Getenv("UPSAVE_METHOD") != "" {
		c.UpSaveMethod = os.Getenv("UPSAVE_METHOD")
	}
	if c.UpSaveMethod == "" {
		c.UpSaveMethod = "POST"
	}
	if os.Getenv("UPSAVE_KEY") != "" {
		c.UpSaveKey = os.Getenv("UPSAVE_KEY")
	}
	if c.UpSaveKey == "" {
		c.UpSaveKey = "userCookie"
	}
	if os.Getenv("DB_ENABLE") == "true" || os.Getenv("DB_ENABLE") == "1" {
		c.DbConf.DbEnable = true
	}
	if os.Getenv("DB_HOST") != "" {
		c.DbConf.DbHost = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		c.DbConf.DbPort = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_USER") != "" {
		c.DbConf.DbUser = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASS") != "" {
		c.DbConf.DbPass = os.Getenv("DB_PASS")
	}
	if os.Getenv("DB_DATABASE") != "" {
		c.DbConf.DbDatabase = os.Getenv("DB_DATABASE")
	}
	if os.Getenv("DB_TYPE") != "" {
		c.DbConf.DbType = os.Getenv("DB_TYPE")
	}
	return c
}

// 保存配置文件
func (c *Conf) Save(p string) error {
	s, _ := json.MarshalIndent(c, "", "\t")
	return util.WriteAllText(p, string(s))
}
