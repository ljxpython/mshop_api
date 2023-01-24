package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"time"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
	//刚才设置的环境变量 想要生效 我们必须得重启goland
}

func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)
	//fmt.Printf("user-web/%s-debug.yaml\n", configFilePrefix)
	fmt.Println(configFileName)
	zap.S().Infof(configFileName)
	if debug {
		configFileName = fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)
	}

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 这个对象如何在其他对象中使用 - 全局变量
	//serverConfig := &config.ServerConifg{}

	//if err := v.Unmarshal(serverConfig); err != nil {
	//	panic(err)
	//}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	//fmt.Println(serverConfig)
	zap.S().Infof("打印此处的日志：")
	zap.S().Infof(global.ServerConfig.UserSrverConfig.Host)
	//zap.S().Infof(string(global.ServerConfig.Port))
	fmt.Println(global.ServerConfig.Name)
	fmt.Println(global.ServerConfig.UserSrverConfig.Host)
	fmt.Println(global.ServerConfig.UserSrverConfig.Port)
	zap.S().Infof("配置信息：%v,%v", global.ServerConfig.UserSrverConfig.Port, global.ServerConfig.UserSrverConfig.Host)
	//fmt.Printf("%V", v.Get("name"))

	//viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file channed： ", e.Name)
		zap.S().Infof("配置文件产生变化： %s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		fmt.Println(global.ServerConfig)
		zap.S().Infof("配置信息：%v", global.ServerConfig)
	})

	time.Sleep(time.Second * 3)

}

func InitConfig1() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: &v", global.NacosConfig)

	//从nacos中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		//RotateTime:          "1h",
		//MaxAge:              3,
		LogLevel: "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		panic(err)
	}
	//fmt.Println(content) //字符串 - yaml
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}
	fmt.Println(&global.ServerConfig)
	//zap.S().Infof(&global.ServerConfig.(string))
}
