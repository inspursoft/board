"use strict";
exports.__esModule = true;
var fs = require("fs");
var child_process = require("child_process");
var CopyFiles = /** @class */ (function () {
  function CopyFiles() {
    this.sourceArr = Array();
    this.sourceArr.push("node_modules/@clr/ui/clr-ui.min.css");
    this.sourceArr.push("node_modules/@clr/icons/clr-icons.min.css");
    this.sourceArr.push("node_modules/echarts/dist/echarts.min.js");
    this.sourceArr.push("node_modules/zone.js/dist/zone.min.js");
    this.sourceArr.push("node_modules/@clr/icons/clr-icons.min.js");
    this.sourceArr.push("src/styles.css");
    this.sourceArr.push("src/favicon.ico");
    this.sourceArr.push("src/rollup/index.html");
  }
  CopyFiles.copyImages = function (source, target) {
    child_process.spawn("cp", ["-r", source, target]);
  };
  CopyFiles.prototype.copyExecute = function () {
    this.sourceArr.forEach(function (value) {
      var pathSource = value.split("/");
      var targetPath = "dist/" + pathSource[pathSource.length - 1];
      fs.createReadStream(value).pipe(fs.createWriteStream(targetPath));
    });
    CopyFiles.copyImages("src/images", "dist/");
  };
  return CopyFiles;
}());
exports.CopyFiles = CopyFiles;
var c = new CopyFiles();
c.copyExecute();
