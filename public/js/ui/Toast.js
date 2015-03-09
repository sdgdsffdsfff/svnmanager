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
			this.element.removeClass('show-toast');
            core.transitionEnd(this.element, function(){
                this.element.remove();
            }.bind(this));
		},
		visible: function(fn){
			this.options.onShow.call(this);
			core.body.append( this.element );
            core.delay(function(){
                this.element.addClass('show-toast');
            }.bind(this), 50);
            core.delay(function(){
                this.hide();
                fn();
            }.bind(this), this.duration)
		}
	},{
		LENGTH_LONG: 5000,
		LENGTH_SHORT: 3000,
		
		TEMPLATE: '<div class="ui-toast"></div>',
		primaryId: 0,
		queue: [],
		processing: false,
        container: null,
		makeText: function(text, duration, options){
            return new Toast(text, duration, options);
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