package template

const  HtmlContent=`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- 上述3个meta标签*必须*放在最前面，任何其他内容都*必须*跟随其后！ -->
    <title>BackEnd Sever</title>

    <!-- Bootstrap -->
    <link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css"
          integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
    <!-- HTML5 shim 和 Respond.js 是为了让 IE8 支持 HTML5 元素和媒体查询（media queries）功能 -->
    <!-- 警告：通过 file:// 协议（就是直接将 html 页面拖拽到浏览器中）访问页面时 Respond.js 不起作用 -->
    <!--[if lt IE 9]>
    <script src="https://cdn.bootcss.com/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="https://cdn.bootcss.com/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->

    <!-- 加载 Bootstrap 的所有 JavaScript 插件。你也可以根据需要只加载单个插件。 -->
    <script src="https://cdn.bootcss.com/jquery/1.12.4/jquery.min.js"></script>
    <script src="https://cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
    <script src="./{{.StructName}}.js"></script>
    <style>
        .table {
            background-color: #fff;
            margin-bottom: 20px;
            margin-top: 20px;
        }

        .table caption {
            position: relative;
            height: 46px;
            line-height: 46px;
            padding-left: 50px;
            padding-right: 40px;
            text-align: left;
            font-size: 21px;
            color: #0088CC;
            border-bottom: 1px solid #ddd;
            margin-bottom: 10px;
        }

        .table caption .btn-group {
            float: right;
        }

        .table td, .table th {
            text-align: center;
            vertical-align: middle;
        }

        select {
            width: 90%;
            margin-bottom: 0;
        }

        button {
            display: block;
            margin: 6px;
        }

    </style>
</head>
<body>
<div class="row">
    <div class="col-md-4">
        <nav aria-label="Page navigation">
            <ul class="pagination">
                <li>
                    <a onclick="GetPage(0)">
                        First
                    </a>
                </li>

                <li>
                    <a onclick="getPrevPage()">
                        Pre
                    </a>
                </li>

                <li>
                    <a onclick="getNextPage()">
                        Next
                    </a>
                </li>

                <li>
                    <a onclick="GetPage(AllPage-1)">
                        Last
                    </a>
                </li>

            </ul>
        </nav>
    </div>
</div>

<table class="table .table-striped" id="test">


</table>

<!-- 添加和编辑问题 -->
<div class="modal fade" id="EditModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span>
                </button>
            </div>
            <div class="modal-body">
                <form class="form-horizontal">
                    {{range $index,$A := .Items }}
                      
					  {{if eq $A.ItemName "Id"}}
 						<fieldset disabled>
                        <div class="form-group">
                            <label for="inputId" class="col-sm-2 control-label">Id:</label>
                            <div class="col-sm-2">
                                <input type="number" class="form-control" id="inputId" placeholder="1" value="0">
                            </div>
                        </div>
                  		  </fieldset>
					 {{else if eq $A.ItemType "int"}}
     							<div class="form-group">
					<label for="input{{$A.ItemName}}" class="col-sm-2 control-label">{{$A.ItemName}}:</label>
    				<input type="number" id="input{{$A.ItemName}}" class="form-control" id="exampleInputEmail1" placeholder="Email">
  					</div>
   					{{else if eq $A.ItemType "array"}}
						
   						//json.{{$A.ItemName}}.push(document.getElementById('').value);
					{{else if eq $A.ItemType "bool"}}
								<div class="form-group">
					<label for="input{{$A.ItemName}}" class="col-sm-2 control-label">{{$A.ItemName}}:</label>
    				<input type="checkbox" id="input{{$A.ItemName}}" class="form-control" id="exampleInputEmail1" placeholder="Email">
  					</div>
					{{else if eq $A.ItemType "float"}}
							<div class="form-group">
					<label for="input{{$A.ItemName}}" class="col-sm-2 control-label">{{$A.ItemName}}:</label>
    				<input type="number" id="input{{$A.ItemName}}" class="form-control" id="exampleInputEmail1" placeholder="Email">
  					</div>
   					{{else }}
					<div class="form-group">
					<label for="input{{$A.ItemName}}" class="col-sm-2 control-label">{{$A.ItemName}}:</label>
    				<input type="text" id="input{{$A.ItemName}}" class="form-control" id="exampleInputEmail1" placeholder="Email">
  					</div>
                    {{end}}
					{{end}}
                </form>
            </div>

            <div class="modal-footer">
                <button type="button" data-dismiss="modal" id="Btn_Submit" class="btn btn-primary ">Submit</button>
                <button type="button" data-dismiss="modal" class="btn btn-danger">Clear</button>
            </div>
        </div>
    </div>
</div>

<!-- Modal -->
<div class="modal fade" id="DeleteModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">
    <div class="modal-dialog" role="document">
        <div class="modal-content">

            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span>
                </button>
                <h4 class="modal-title">Alert</h4>
            </div>
            <div class="modal-body">
                <p>Delete？</p>
                <h4 id="delete{{.StructName}}"/>
            </div>
            <div class="modal-footer">
                <button type="button" data-dismiss="modal" class="btn btn-primary" id="Btn_Delete">Delete</button>
            </div>
        </div>
    </div>
</div>

<div class="row">
    <div class="col-md-4">
        <button type="button" class="btn btn-primary btn-lg" onclick="Add{{.StructName}}()">
            添加问题
        </button>
    </div>
    <div class="col-md-4">
        <nav aria-label="Page navigation">
            <ul class="pagination">
                <li>
                    <a onclick="GetPage(0)">
                        First
                    </a>
                </li>

                <li>
                    <a onclick="getPrevPage()">
                        Pre
                    </a>
                </li>

                <li>
                    <a onclick="getNextPage()">
                        Next
                    </a>
                </li>

                <li>
                    <a onclick="GetPage(AllPage-1)">
                        Last
                    </a>
                </li>

            </ul>
        </nav>
    </div>
</div>

<script>

    getData(0);
    var Btn_Submit = document.getElementById('Btn_Submit');
    Btn_Submit.onclick = function () {
        addContent(getFormData());
    }
</script>
</body>
</html>`