function getPrevPage() {
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
    $.post("/test/updateinfo", content, function (data, status) {
        if (status == 'success') {
            getData(CurrentPage)
        } else {
            alert("Submit fail")
        }
    })
}

var CurrentPage = 0;
var infos;
var AllPage;

//  getData(CurrentPage);

function getData(currentPage) {
    $.get("/test/infos?page=" + currentPage, function (data, status) {
        if (status == 'success') {
            var self = $("#test");
            self.empty();
            addHead();
            var obj = jQuery.parseJSON(data);
            infos = obj.infos;
            AllPage = obj.AllPage;
            for (x in Questions) {
                console.log(x)
                addContent(obj.Questions[x], x)
            }
        }
    })
}

function confirmDelete(index) {
    var json = Questions[index];
    document.getElementById("deleteinfo").innerHTML = json.info;
    $('#DeleteModal').modal('show')
    document.getElementById("Btn_Delete").onclick = function () {
        deleteinfo(json.Id)
    }
}

function deleteinfo(index) {
    $.post("/test/deleteinfo?id=" + index, function (data, status) {
        if (status == 'success') {
            getData(CurrentPage)
        } else {
            alert("Delete fail")
        }
    })

}


function getFormData() {
    var json = {


    'Id': 0,



    'Name': '',


    };

   
   
    json.Id = parseInt(document.getElementById('inputId').value);
   
   
   
   json.Name = document.getElementById('inputName').value;
   
   

    console.log(JSON.stringify(json))
    updateOrCreateData(json)
    return json;
}

function addHead() {
    var str = '   <thead>\n' +
        '    <tr>\n' +
         
         '        <th>Id</th>\n' +
         
         '        <th>Name</th>\n' +
         
        '    </tr>\n' +
        '    </thead>'

    var self = $("#test");
    self.append(str)
}

function addContent(json, x) {
    var self = $("#test");
    var $tr = '';

    $tr += '<tr class="active" id=content-' + json.Id +
        '><td scope="row">' + json.Id + '</td>';
     
             $tr += '<td class="active" width=100px>' + json.Id + '</td>';
     
             $tr += '<td class="active" width=100px>' + json.Name + '</td>';
     
    $tr += '<td class="active"><button type="button" class="btn" onclick= "update(' + x + ')"    padding-left=50px>Update</button>' +
        '<button type="button" class="btn  btn-warning"   onclick= "confirmDelete(' + x + ')">Delete</button></td></tr>';
    self.append($tr);
}

function update(index) {
    var json = Questions[index];
    
          document.getElementById("inputId").value = json.Id;
    
          document.getElementById("inputName").value = json.Name;
    
    $('#EditModal').modal('show')
}

function Addinfo() {

 
    document.getElementById("inputId").value = 0;
   

 
   document.getElementById("inputName").value = "";
   

    $('#EditModal').modal('show')
}
