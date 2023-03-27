package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"openai/types"
	"openai/utils"
	"strconv"
)

// ConfigSetHandle set configs
func (s *Server) ConfigSetHandle(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Model
	if model, ok := data["model"]; ok {
		s.Config.Chat.Model = model
	}

	// Temperature
	if temperature, ok := data["temperature"]; ok {
		v, err := strconv.ParseFloat(temperature, 32)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{
				Code:    types.InvalidParams,
				Message: "temperature must be a float parameter",
			})
			return
		}
		s.Config.Chat.Temperature = float32(v)
	}

	// max_tokens
	if maxTokens, ok := data["max_tokens"]; ok {
		v, err := strconv.Atoi(maxTokens)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{
				Code:    types.InvalidParams,
				Message: "max_tokens must be a int parameter",
			})
			return
		}
		s.Config.Chat.MaxTokens = v
	}

	// enable Context
	if enableContext, ok := data["enable_context"]; ok {
		v, err := strconv.ParseBool(enableContext)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{
				Code:    types.InvalidParams,
				Message: "enable_context must be a bool parameter",
			})
			return
		}
		s.Config.Chat.EnableContext = v
	}

	// enable auth
	if enableAuth, ok := data["enable_auth"]; ok {
		v, err := strconv.ParseBool(enableAuth)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{
				Code:    types.InvalidParams,
				Message: "enable_auth must be a bool parameter",
			})
			return
		}
		s.Config.EnableAuth = v
	}

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg})
}

// AddToken 添加 Token
func (s *Server) AddToken(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	var name = data["name"]
	var maxCalls = data["max_calls"]
	if name == "" || maxCalls == "" {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	}

	n, err := strconv.Atoi(maxCalls)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{
			Code:    types.InvalidParams,
			Message: "enable_auth must be a int parameter",
		})
		return
	}

	var tokens = GetTokens()
	if utils.ContainToken(tokens, name) {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Token " + name + " already exists"})
		return
	}

	err = PutToken(types.Token{Name: name, MaxCalls: n, RemainingCalls: n})
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save configs"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: GetTokens()})
}

// RemoveToken 删除 Token
func (s *Server) RemoveToken(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if token, ok := data["token"]; ok {
		err = RemoveToken(token)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save configs"})
			return
		}
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: GetTokens()})
}

// AddApiKey 添加一个 API key
func (s *Server) AddApiKey(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	if key, ok := data["api_key"]; ok && len(key) > 20 {
		s.Config.Chat.ApiKeys = append(s.Config.Chat.ApiKeys, key)
	}

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.Chat.ApiKeys})
}

// RemoveApiKey 移除一个 API key
func (s *Server) RemoveApiKey(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if key, ok := data["api_key"]; ok {
		for i, v := range s.Config.Chat.ApiKeys {
			if v == key {
				s.Config.Chat.ApiKeys = append(s.Config.Chat.ApiKeys[:i], s.Config.Chat.ApiKeys[i+1:]...)
			}
		}
	}

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.Chat.ApiKeys})
}

// ListApiKeys 获取 API key 列表
func (s *Server) ListApiKeys(c *gin.Context) {
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.Chat.ApiKeys})
}

func (s *Server) GetChatRoles(c *gin.Context) {
	//var rolesOrder = []string{"gpt", "programmer", "teacher", "artist", "philosopher", "lu-xun", "english_trainer", "seller"}
	//var roles = make([]interface{}, 0)
	//for _, k := range rolesOrder {
	//	if v, ok := s.Config.ChatRoles[k]; ok && v.Enable {
	//		roles = append(roles, struct {
	//			Key  string `json:"key"`
	//			Name string `json:"name"`
	//			Icon string `json:"icon"`
	//		}{
	//			Key:  v.Key,
	//			Name: v.Name,
	//			Icon: v.Icon,
	//		})
	//	}
	//}
	//c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: roles})
}

// UpdateChatRole 更新某个聊天角色信息，这里只允许更改名称以及启用和禁用角色操作
func (s *Server) UpdateChatRole(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	key := data["key"]
	if key == "" {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Must specified the role key"})
		return
	}

	//role := s.Config.ChatRoles[key]
	//if enable, ok := data["enable"]; ok {
	//	v, err := strconv.ParseBool(enable)
	//	if err != nil {
	//		c.JSON(http.StatusOK, types.BizVo{
	//			Code:    types.InvalidParams,
	//			Message: "enable must be a bool parameter",
	//		})
	//		return
	//	}
	//	role.Enable = v
	//}

	//if name, ok := data["name"]; ok {
	//	role.Name = name
	//}
	//if helloMsg, ok := data["hello_msg"]; ok {
	//	role.HelloMsg = helloMsg
	//}
	//if icon, ok := data["icon"]; ok {
	//	role.Icon = icon
	//}
	//
	//s.Config.ChatRoles[key] = role
	//
	//// 保存配置文件
	//err = types.SaveConfig(s.Config, s.ConfigPath)
	//if err != nil {
	//	c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg})
}

// AddProxy 添加一个代理
func (s *Server) AddProxy(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if proxy, ok := data["proxy"]; ok {
		if !utils.ContainsStr(s.Config.ProxyURL, proxy) {
			s.Config.ProxyURL = append(s.Config.ProxyURL, proxy)
		}
	}

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.ProxyURL})
}

func (s *Server) RemoveProxy(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if proxy, ok := data["proxy"]; ok {
		for i, v := range s.Config.ProxyURL {
			if v == proxy {
				s.Config.ProxyURL = append(s.Config.ProxyURL[:i], s.Config.ProxyURL[i+1:]...)
				break
			}
		}
	}

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.ProxyURL})
}
