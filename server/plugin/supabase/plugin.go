package supabase

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/plugin/supabase/initialize"
	interfaces "github.com/icosmos-space/iadmin/server/utils/plugin/v2"
)

var _ interfaces.Plugin = (*plugin)(nil)

var Plugin = new(plugin)

type plugin struct{}

func init() {
	interfaces.Register(Plugin)
}

func (p *plugin) Register(group *gin.Engine) {
	ctx := context.Background()
	initialize.Viper()
	initialize.Api(ctx)
	initialize.Menu(ctx)
	initialize.Dictionary(ctx)
	initialize.Gorm(ctx)
	initialize.Router(group)
}
