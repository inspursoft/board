import os
import sys

def main():
    os.system('cd ..')
    tag = os.popen('git describe --tags').readline()
    print (tag)
    tagDir = tag.split("\n")[0]
    resultDir = sys.argv[1] +  "/" + tagDir
    resultHtml = resultDir + "/index.html"
    if not os.path.exists(resultDir):
        os.makedirs(resultDir)
    inf = open(resultHtml,'w')
    inf.write("<html>")
    inf.write("<body>")
    inf.write("<table>")

    f = open("out.temp")
    lines = f.readlines()
    for line in lines:
        inf.write("<tr>\n")
        inf.write("<td>" + line + "</td>\n")
        inf.write("</tr>\n")
    inf.write("</table>\n")
    inf.write("</body>\n")
    inf.write("</html>")
    f.close()
    inf.close()
    
    cov = os.popen("cat out.temp|grep \"total\"|awk '{print $NF}'").readline()
    covDir = resultDir + "/index"
    covHtml = covDir + "/index.html"
    if not os.path.exists(covDir):
        os.makedirs(covDir)
    covf = open(covHtml,"w")
    covf.write('<html>')
    covf.write('<body>')
    covf.write(cov)
    covf.write('</html>')
    covf.write('</body>')
    covf.close()

    tagDir = resultDir + "/tag"
    tagHtml = tagDir + "/index.html"
    if not os.path.exists(tagDir):
        os.makedirs(tagDir)
    tagf = open(tagHtml, "w")
    tagf.write("<html>")
    tagf.write("<body>")
    tagf.write(tag)
    tagf.write("</html>")
    tagf.write("</body>")
    tagf.close()

if __name__ == "__main__":
    main()
