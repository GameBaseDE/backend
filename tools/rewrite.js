fs = require('fs');
process.chdir('tmp');
module.paths.push('node_modules');
gofmt = require('gofmt.js');

// match all go funcs like:
// func DeleteContainer(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{})
// }
// and replace them with:
// func DeleteContainer(c *gin.Context) {
// 	gamebase.DeleteContainer(c)
// }
const endpointHandler = /^func ([a-zA-Z]+)\(c \*gin\.Context\) {\n\t[a-zA-Z(){}., ]*\n}$/gm;

const REWRITE = {
    'api_auth.go': content =>
        content
            .replace('"net/http"', '')
            .replace(endpointHandler, 'func $1(c *gin.Context) {\n\t//TODO: authenticator.$1(c)\n}'),
    'api_gameserver.go': content =>
        content
            .replace('"net/http"', '')
            .replace(endpointHandler, 'func $1(c *gin.Context) {\n\tauthenticator.$1(c)\n}'),
    'routers.go': content => content // add cors header
            .replace('router := gin.Default()', 'router := gin.Default()\n\t\n	corsConfig := cors.DefaultConfig()\n	corsConfig.AllowAllOrigins = true\n	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers")\n	router.Use(cors.New(corsConfig))')
        .replace('"github.com/gin-gonic/gin"', '"github.com/gin-contrib/cors"\n\t"github.com/gin-gonic/gin"')
        + '\n\nvar api = NewAPI()\nvar authenticator = newHttpRequestAuthenticator()',
};

process.chdir('out/go');
for (const file of fs.readdirSync(".")) {
    const f = REWRITE[file];
    if (f) {
        console.log(`Replacing file out/go/${file}`)
        let content = fs.readFileSync(file, 'utf-8');
        content = gofmt(f(content));
        fs.writeFileSync(file, content);
    } else {
        console.log(`Skipping file out/go/${file}`)
    }
}