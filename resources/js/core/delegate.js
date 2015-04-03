//事件代理调用，为页面上所有拥有action-type的a标签和button绑定点击事件。
//需要在初始化时调用
define(['kernel'],
function(core){
	return function(){
		$('body').delegate('a[action-type], button[action-type]','click',function(e){
			var anchor = $(this), type = this.getAttribute('action-type'), action = null, isInline = type.indexOf('@') == 0;
			
			if( isInline ){
				type = type.substring(1, type.length);
			}
			
			if(core.actions && (action = core.actions[type]) ){
				action.call(anchor, e)
			}else{
				//使用圈号开头的事件将会等到内部代码注册事件
				if( isInline ){
					core.console.log('waiting inline action: ' + type + '...')
				}else{
				//否则使用require加载事件并执行
					require(['mo/actions/'+type], function(action){
						try{
							action.call(anchor)
						}catch(e){
							core.console.log('action-Error:'+e)
						}
					})
				}
			}
			return false;
		});
	};
})