// THIS FILE IS AUTO-GENERATED FROM home.js
// DO NOT EDIT.
package ui
func init() {
	Files.Add("home.js", "/*\nCopyright 2011 Google Inc.\n\nLicensed under the Apache License, Version 2.0 (the \"License\");\nyou may not use this file except in compliance with the License.\nYou may obtain a copy of the License at\n\n     http://www.apache.org/licenses/LICENSE-2.0\n\nUnless required by applicable law or agreed to in writing, software\ndistributed under the License is distributed on an \"AS IS\" BASIS,\nWITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\nSee the License for the specific language governing permissions and\nlimitations under the License.\n*/\n\n// CamliHome namespace to contain the global vars\nvar CamliHome = {};\n\nfunction btnCreateNewPermanode(e) {\n    camliCreateNewPermanode(\n        {\n            success: function(blobref) {\n               window.location = \"./?p=\" + blobref;\n            },\n            fail: function(msg) {\n                alert(\"create permanode failed: \" + msg);\n            }\n        });\n}\n\nfunction handleFormSearch(e) {\n    e.stopPropagation();\n    e.preventDefault();\n\n    var input = document.getElementById(\"inputSearch\");\n    var btn = document.getElementById(\"btnSearch\");\n\n    if (input.value == \"\") {\n        return;\n    }\n\n    var query = input.value.split(/\\s*,\\s*/);\n    window.location = \"./search.html?q=\" + query[0] + \"&t=tag\";\n}\n\nfunction indexOnLoad(e) {\n   var btnNew = document.getElementById(\"btnNew\");\n    if (!btnNew) {\n        alert(\"missing btnNew\");\n    }\n    btnNew.addEventListener(\"click\", btnCreateNewPermanode);\n    camliGetRecentlyUpdatedPermanodes({ success: indexBuildRecentlyUpdatedPermanodes });\n    formSearch.addEventListener(\"submit\", handleFormSearch);\n\n    if (disco && disco.uploadHelper) {\n        var uploadForm = document.getElementById(\"uploadform\");\n        uploadform.action = disco.uploadHelper;\n        document.getElementById(\"fileinput\").disabled = false;\n        document.getElementById(\"filesubmit\").disabled = false;\n        var chkRollSum = document.getElementById(\"chkrollsum\");\n        chkRollSum.addEventListener(\"change\", function (e) {\n                                        if (chkRollSum.checked) {\n                                            if (disco.uploadHelper.indexOf(\"?\") == -1) {\n                                                uploadform.action = disco.uploadHelper + \"?rollsum=1\";\n                                            } else {\n                                                uploadform.action = disco.uploadHelper + \"&rollsum=1\";\n                                            }\n                                        } else {\n                                            uploadform.action = disco.uploadHelper;\n                                        }\n                                    });\n    }\n}\n\nfunction indexBuildRecentlyUpdatedPermanodes(searchRes) {\n    var div = document.getElementById(\"recent\");\n    div.innerHTML = \"\";\n    for (var i = 0; i < searchRes.recent.length; i++) {\n        var result = searchRes.recent[i];      \n        var pdiv = document.createElement(\"li\");\n        var alink = document.createElement(\"a\");\n        alink.href = \"./?p=\" + result.blobref;\n        alink.innerText = camliBlobTitle(result.blobref, searchRes);\n        pdiv.appendChild(alink);\n        div.appendChild(pdiv);\n    }\n}\n\nwindow.addEventListener(\"load\", indexOnLoad);\n");
}
