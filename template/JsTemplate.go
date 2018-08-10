package template

const JSContent  =`function getPrevPage() {
    if (CurrentPage > 0) {
        CurrentPage--;
        getData(CurrentPage);
    }
}

function GetPage(index) {
    CurrentPage = index;
    getData(CurrentPage);
}


function getNextPage() {
    CurrentPage++;
    getData(CurrentPage);
}


function updateOrCreateData(json) {
    var content = {
        'key': ''
    }
    content.key = JSON.stringify(json)
    $.post("/{{.ModuleName}}/update{{.StructName}}", content, function (data, status) {
        if (status == 'success') {
            getData(CurrentPage)
        } else {
            alert("Submit fail")
        }
    })
}

var CurrentPage = 0;
var {{.StructName}}s;
var AllPage;

//  getData(CurrentPage);

function getData(currentPage) {
    $.get("/{{.ModuleName}}/get{{.StructName}}s?page=" + currentPage, function (data, status) {
        if (status == 'success') {
            var self = $("#test");
            self.empty();
            addHead();
            var obj = jQuery.parseJSON(data);
            {{.StructName}}s = obj.{{.StructName}}s;
            AllPage = obj.AllPage;
            for (x in {{.StructName}}s) {
                console.log(x)
                addContent(obj.{{.StructName}}s[x], x)
            }
        }
    })
}

function confirmDelete(index) {
    var json = Questions[index];
    document.getElementById("delete{{.StructName}}").innerHTML = json.{{.StructName}};
    $('#DeleteModal').modal('show')
    document.getElementById("Btn_Delete").onclick = function () {
        delete{{.StructName}}(json.Id)
    }
}

function delete{{.StructName}}(index) {
    $.post("/{{.ModuleName}}/delete{{.StructName}}?id=" + index, function (data, status) {
        if (status == 'success') {
            getData(CurrentPage)
        } else {
            alert("Delete fail")
        }
    })

}


function getFormData() {
    var json = {
{{range $index,$A := .Items }}
{{if eq $A.ItemType "int"}} '{{$A.ItemName}}': 0,
{{else if  eq $A.ItemType "array"}} '{{$A.ItemName}}': [],
{{else if eq $A.ItemType "float"}} '{{$A.ItemName}}': 0,
{{else if eq $A.ItemType "string"}} '{{$A.ItemName}}': '',
{{else}}'{{$A.ItemName}}': '',
{{end}}{{end}}
    };
   {{range $index,$A := .Items }}
   {{if eq $A.ItemType "int"}} json.{{$A.ItemName}} = parseInt(document.getElementById('input{{$A.ItemName}}').value);
   {{else if eq $A.ItemType "array"}}//json.{{$A.ItemName}}.push(document.getElementById('').value);
   {{else if eq $A.ItemType "float"}}json.{{$A.ItemName}} =parseFloat( document.getElementById('input{{$A.ItemName}}').value);
   {{else}}json.{{$A.ItemName}} = document.getElementById('input{{$A.ItemName}}').value;
   {{end}}{{end}}
    console.log(JSON.stringify(json))
    updateOrCreateData(json)
    return json;
}

function addHead() {
    var str = '   <thead>\n' +
        '    <tr>\n' +
         {{range $index,$A := .Items }}
         '        <th>{{$A.ItemName}}</th>\n' +
         {{end}}
        '    </tr>\n' +
        '    </thead>'
    var self = $("#test");
    self.append(str)
}

function addContent(json, x) {
    var self = $("#test");
    var $tr = '';
     {{range $index,$A := .Items }}
			{{if eq $A.ItemName "Id"}}
 			$tr += '<tr class="active" id=content-' + json.Id +
        		'><td scope="row">' + json.Id + '</td>';
			{{else}}
				$tr += '<td class="active" width=100px>' + json.{{$A.ItemName}} + '</td>';
			{{end}}
     {{end}}
    $tr += '<td class="active"><button type="button" class="btn" onclick= "update(' + x + ')"    padding-left=50px>Update</button>' +
        '<button type="button" class="btn  btn-warning"   onclick= "confirmDelete(' + x + ')">Delete</button></td></tr>';
    self.append($tr);
}

function update(index) {
    var json = {{.StructName}}s[index];
    {{range $index,$A := .Items }}
          document.getElementById("input{{$A.ItemName}}").value = json.{{$A.ItemName}};
    {{end}}
    $('#EditModal').modal('show')
}

function Add{{.StructName}}() {
{{range $index,$A := .Items }}
 {{if eq $A.ItemType "int"}}
    document.getElementById("input{{$A.ItemName}}").value = 0;
   {{else if eq $A.ItemType "array"}}
   //document.getElementById('input{{$A.ItemName}}0').value = "";
   {{else if eq $A.ItemType "string"}}
   document.getElementById("input{{$A.ItemName}}").value = "";
   {{else if eq $A.ItemType "float"}}
	document.getElementById("input{{$A.ItemName}}").value = 0;
   {{else }}
   document.getElementById("input{{$A.ItemName}}").value = "";
   {{end}}
{{end}}
    $('#EditModal').modal('show')
}
`
