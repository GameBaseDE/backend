fs = require('fs');
process.chdir('tmp');
module.paths.push('node_modules');
gofmt = require('gofmt.js');

const REWRITE = {
    'api_api.go': content => {
        // match all go funcs like:
        // func DeleteContainer(c *gin.Context) {
        // 	c.JSON(http.StatusOK, gin.H{})
        // }
        // and replace them with:
        // func DeleteContainer(c *gin.Context) {
        // 	gamebase.DeleteContainer(c)
        // }
        const regExp = /^func ([a-zA-Z]+)\(c \*gin\.Context\) {\n\t[a-zA-Z(){}., ]*\n}$/gm;
        return content
            .replace('"net/http"', '')
            .replace(regExp, 'func $1(c *gin.Context) {\n\t$1_(c)\n}');
    },
    'routers.go': content => content // add cors header
        .replace('router := gin.Default()', 'router := gin.Default()\n\trouter.Use(cors.Default())')
        .replace('"github.com/gin-gonic/gin"', '"github.com/gin-contrib/cors"\n\t"github.com/gin-gonic/gin"'),
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