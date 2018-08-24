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
    $.post("/User/updateBaseUser", content, function (data, status) {
        if (status == 'success') {
            getData(CurrentPage)
        } else {
            alert("Submit fail")
        }
    })
}

var CurrentPage = 0;
var BaseUsers;
var AllPage;

//  getData(CurrentPage);

function getData(currentPage) {
    $.get("/User/getBaseUsers?page=" + currentPage, function (data, status) {
        if (status == 'success') {
            var self = $("#test");
            self.empty();
            addHead();
            var obj = jQuery.parseJSON(data);
            BaseUsers = obj.BaseUsers;
            AllPage = obj.AllPage;
            for (x in BaseUsers) {
                console.log(x)
                addContent(obj.BaseUsers[x], x)
            }
        }
    })
}

function confirmDelete(index) {
    var json = Questions[index];
    document.getElementById("deleteBaseUser").innerHTML = json.BaseUser;
    $('#DeleteModal').modal('show')
    document.getElementById("Btn_Delete").onclick = function () {
        deleteBaseUser(json.Id)
    }
}

function deleteBaseUser(index) {
    $.post("/User/deleteBaseUser?id=" + index, function (data, status) {
        if (status == 'success') {
            getData(CurrentPage)
        } else {
            alert("Delete fail")
        }
    })

}


function getFormData() {
    var json = {

 'Name': '',

 'Phone': '',

 'Password': '',

 'UserID': '',

 'Id': 0,

 'CreatTime': 0,

    };
   
   json.Name = document.getElementById('inputName').value;
   
   json.Phone = document.getElementById('inputPhone').value;
   
   json.Password = document.getElementById('inputPassword').value;
   
   json.UserID = document.getElementById('inputUserID').value;
   
    json.Id = parseInt(document.getElementById('inputId').value);
   
    json.CreatTime = parseInt(document.getElementById('inputCreatTime').value);
   
    console.log(JSON.stringify(json))
    updateOrCreateData(json)
    return json;
}

function addHead() {
    var str = '   <thead>\n' +
        '    <tr>\n' +
         
         '        <th>Name</th>\n' +
         
         '        <th>Phone</th>\n' +
         
         '        <th>Password</th>\n' +
         
         '        <th>UserID</th>\n' +
         
         '        <th>Id</th>\n' +
         
         '        <th>CreatTime</th>\n' +
         
        '    </tr>\n' +
        '    </thead>'
    var self = $("#test");
    self.append(str)
}

function addContent(json, x) {
    var self = $("#test");
    var $tr = '';
     
			
				$tr += '<td class="active" width=100px>' + json.Name + '</td>';
			
     
			
				$tr += '<td class="active" width=100px>' + json.Phone + '</td>';
			
     
			
				$tr += '<td class="active" width=100px>' + json.Password + '</td>';
			
     
			
				$tr += '<td class="active" width=100px>' + json.UserID + '</td>';
			
     
			
 			$tr += '<tr class="active" id=content-' + json.Id +
        		'><td scope="row">' + json.Id + '</td>';
			
     
			
				$tr += '<td class="active" width=100px>' + json.CreatTime + '</td>';
			
     
    $tr += '<td class="active"><button type="button" class="btn" onclick= "update(' + x + ')"    padding-left=50px>Update</button>' +
        '<button type="button" class="btn  btn-warning"   onclick= "confirmDelete(' + x + ')">Delete</button></td></tr>';
    self.append($tr);
}

function update(index) {
    var json = BaseUsers[index];
    
          document.getElementById("inputName").value = json.Name;
    
          document.getElementById("inputPhone").value = json.Phone;
    
          document.getElementById("inputPassword").value = json.Password;
    
          document.getElementById("inputUserID").value = json.UserID;
    
          document.getElementById("inputId").value = json.Id;
    
          document.getElementById("inputCreatTime").value = json.CreatTime;
    
    $('#EditModal').modal('show')
}

function AddBaseUser() {

 
   document.getElementById("inputName").value = "";
   

 
   document.getElementById("inputPhone").value = "";
   

 
   document.getElementById("inputPassword").value = "";
   

 
   document.getElementById("inputUserID").value = "";
   

 
    document.getElementById("inputId").value = 0;
   

 
    document.getElementById("inputCreatTime").value = 0;
   

    $('#EditModal').modal('show')
}
