// Package custom 注册用户自定义插件于此
package custom

import (
	_ "github.com/FloatTech/ZeroBot-Plugin/custom/plugin/antirecall"        // 反撤回
	_ "github.com/FloatTech/ZeroBot-Plugin/custom/plugin/deepseek"        // deepseek
	_ "github.com/FloatTech/ZeroBot-Plugin/custom/plugin/delreply"        // 回复撤回
	_ "github.com/FloatTech/ZeroBot-Plugin/custom/plugin/moehu"        // 图片
)