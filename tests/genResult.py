import os
import sys

def main():
    resultDir = sys.argv[1]
    resultHtml = resultDir + "/index.html"
    if not os.path.exists(resultDir):
        os.makedirs(resultDir)
    print "xxx"
    inf = open(resultHtml,'w')
    inf.write("<html>")
    inf.write("<body>")
    inf.write("<table>")

    f = open("out.temp")
    lines = f.readlines()
    for line in lines:
        print "==========="
        print ("<tr>" + line + "</tr>")
        print ("xxxxxxxxxxxxxxxxxxx")
        inf.write("<tr>\n")
        inf.write("<td>" + line + "</td>\n")
        inf.write("</tr>\n")
    inf.write("</table>\n")
    inf.write("</body>\n")
    inf.write("</html>")
    f.close()
    inf.close()

if __name__ == "__main__":
    main()
