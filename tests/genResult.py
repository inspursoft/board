#!/bin/python

import os
import sys
import shutil

def main():
    os.system('cd ..')
    tag = os.popen('git describe --tags').readline()
    print (tag)
    tagDir = tag.split("\n")[0]
    resultDir = sys.argv[1] +  "/" + tagDir
    resultDir = sys.argv[1] +  "/total" 
    resultHtml = resultDir + "/index.html"
    if not os.path.exists(resultDir):
        os.makedirs(resultDir)
    inf = open(resultHtml,'w')
    inf.write("<html>")
    inf.write("<body>")
    inf.write("<table>")

    #cov = os.popen("cat out.temp|grep \"total\"|awk '{print $NF}'").readline()
    cov = sys.argv[2]
    covDir = resultDir + "/index"
    covHtml = covDir + "/index.html"
    genCovHtml(cov, covDir, covHtml)

    uicov = sys.argv[3] 
    uiCovDir =os.path.join(resultDir, 'coverage', 'cov')
    uiCovHtml = os.path.join(uiCovDir, 'index.html')
    genCovHtml(uicov, uiCovDir, uiCovHtml)


    tagDir = resultDir + "/tag"
    tagHtml = tagDir + "/index.html"
    genCovHtml(tag, tagDir, tagHtml)


    shutil.copy("profile.html", resultHtml)
    print (resultHtml) 

def genCovHtml(cov, covDir, covHtml):
    if not os.path.exists(covDir):
        os.makedirs(covDir)
    covf = open(covHtml,"w")
    covf.write('<html>')
    covf.write('<body>')
    covf.write(cov + "%")
    covf.write('</html>')
    covf.write('</body>')
    covf.close()


if __name__ == "__main__":
    main()
