"use strict";
exports.__esModule = true;
var fs = require("fs");
var CreateMainJs = /** @class */ (function () {
    function CreateMainJs() {
        this.jsContent = Array();
        this.jsContent.push("import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';\n");
        this.jsContent.push("import { enableProdMode } from '@angular/core';\n");
        this.jsContent.push("import { environment } from './environments/environment';\n");
        this.jsContent.push("import { AppModuleNgFactory } from './app/app.module.ngfactory';\n");
        this.jsContent.push("if (environment.production) {\n");
        this.jsContent.push("enableProdMode();\n");
        this.jsContent.push("}\n");
        this.jsContent.push("platformBrowserDynamic().bootstrapModuleFactory(AppModuleNgFactory);\n");
    }
    CreateMainJs.prototype.createMainAotJsFile = function () {
        var _this = this;
        fs.open("out-ngc/src/main.js", "w", function (err, fd) {
            if (err)
                throw err;
            _this.jsContent.forEach(function (value) {
                fs.writeSync(fd, Buffer.from(value));
            });
        });
    };
    return CreateMainJs;
}());
exports.CreateMainJs = CreateMainJs;
var c = new CreateMainJs();
c.createMainAotJsFile();
