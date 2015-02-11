define([
'kernel'
],
function(core){
	var Toast = core.Class.extend({
		
		primaryId: 0,
		winEventName: '',
		enabled: false,
		
		init: function(text, duration, options){

			this.text = text;
			this.duration = duration >= 0 ? duration : Toast.LENGTH_SHORT;
			this.options = $.extend({
				onCreate: $.noop,
				onShow: $.noop,
				onHide: $.noop,
				classStyle: '',
				offset: {
					bottom: 100
				}
			}, options);
			
			this.primaryId = Toast.primaryId++;
			this.create();
		},
		create: function(){
			this.element = $(Toast.TEMPLATE);
			this.element.addClass(this.options.classStyle).html( this.text );
			this.options.onCreate.call(this);
		},
		show: function(){
			Toast.queue.push(this);
			this.enabled = true;
			if( !Toast.processing ){
				Toast.show();
			}
		},
		hide: function(){
			this.options.onHide.call(this);
			core.xresize.unbind(this.winEventName);
			this.element.fadeOut('slow', function(){
				$(this).remove();
			});
		},
		visible: function(fn){
			var self = this;
			
			this.options.onShow.call(this);
			core.body.append(this.element);
			this.winEventName = core.positionFixed(
				this.element,
				{
					centerHorizontal: true,
					offset: this.options.offset,
					onPosition: function(ele){
						ele.fadeIn();
					}
				}
			);
			
			core.delay(
				$.proxy(function(){
					this.hide();
					fn();
				}, this),
				this.duration
			)
		}
	},{
		LENGTH_LONG: 5000,
		LENGTH_SHORT: 3000,
		
		TEMPLATE: '<div class="ui-toast"></div>',
		primaryId: 0,
		queue: [],
		processing: false,
		makeText: function(text, duration, options){
			if( text ){
				return new Toast(text, duration, options);
			}
			return null;
		},
		next: function(){
			
		},
		show: function(){
			
			if( !Toast.queue.length ){
				Toast.processing = false;
				return;
			}
			
			Toast.processing = true;
			
			var toast = Toast.queue.shift();
			
			if( toast ){
				if( toast.enabled ){
					toast.visible( Toast.show )
				}else{
					Toast.show();
				}
			}
		}
	});
	
	return Toast;
});