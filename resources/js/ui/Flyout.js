define(['kernel', 'react'],
function(core, React){
	
	var $window = $(window), $document = $(document);
	
	var Flyout = core.Class.extend({

		init: function (element, options, props) {
        	
        	var flyoutModel;

        	if( typeof element == 'string' ){
    			//use template
    			if(flyoutModel = Flyout.template[ element ]){
    				element = $(flyoutModel.element);
    				options = $.extend({},flyoutModel.options, options)
    			//create element
    			}else{
    				element = $(element)
    			}
    			core.body.append( element )
    		}

            if( element.attr('id') && !document.getElementById(element.attr('id')) ) {
                core.body.append( element );
            }

            this._reactElement = element.data('reactElement');

			if (typeof props !== 'undefined') {
				$.extend(this, props);
			}
    		
            if(element.data('flyout')){ 
               return element.data('flyout');
            }
            
            element.data('flyout', this );
            
            this._baseFlyoutConstructor(element, options);
            
            this.primaryId = Flyout.primaryId++;
            
            Flyout.cache[ this.primaryId ] = this;
            
        },
        _element: null,
        _anchor: null,
        _where: {},
        _event_hide_id: null,
		_isHide: true,
		_currentAnchor: null,
		_currentAnchorBox: {},
        _reactElement: null,
        
        //用于扩展属性
        parasitic: {},

        getRef: function( name ){
            if( typeof name === 'undefined' ){
                return this._reactElement.refs
            }
            if( name in this ){
                return this[name];
            }
            if( name in this._reactElement.refs ){
                return this._reactElement.refs[name];
            }
            return null;
        },

        getReact: function(){
            return this._reactElement;
        },

		_baseFlyoutConstructor: function (element, options, methods) {
        	
        	this.setting = $.extend({
        		onRendered: $.noop,
        		onShow: $.noop,
        		onHide: $.noop,
        		//隐藏时是否也清楚对象
        		destroy: false,
        		//不会被全局隐藏，不会隐藏全局，独立存在
        		alone: false,
        		//点击同一个锚点时是否隐藏
        		sameAnchorHide: true,
        		//隐藏的方向
        		hideDirect: 'default',
        		//自动隐藏的停留时间
        		stayTime: 0,
        		//锚点悬停时是否持续显示
        		anchorStay: false,
        		//隐藏位置偏移量
        		hideOffset: 10,
        		//显示位置偏移量
        		showOffset: 10,
        		//动画时间
        		showDuration: 150,
        		easing: 'swing',
        		classStyle: '',
        		alignment: 'center',
        		placement: 'top',
				arrow: true,
				stop: false
        	}, options);

			$.isPlainObject(methods) && $.extend(this, methods);

        	//元素节点，可修改
        	this.element = element;
        	
        	this.alone = this.setting.alone;
        	
        	//保存原始节点
        	this._originalHTML = element.html();
        	
        	this.element.addClass('ui-flyout').css({
        		position: 'absolute'
        	});

			if (this.setting.stop) {
        		this.element.click(function(e){
        			e.stopPropagation()
        		})
        	}
    		if(this.setting.anchor){
    			this._anchor = options.anchor;
    		}
			this.element.addClass(this.setting.classStyle);
        	
        	//相对于冒点的位置[top, bottom,left, right ]
        	//TODO 这里没有实现auto方法
        	this._placement = this.setting.placement;
        	//对齐方式[left, center, right ]，仅针对placement 为 top 或者 bottom
        	this._alignment = this.setting.alignment;
        	
        	this.setting.onRendered.call( this, element );

			if (this.setting.arrow) {
				this._arrow = $('<span class="flyout-arrow"></span>');
				this.element.append(this._arrow)
			}
        },
		arrow: function () {
			"use strict";
			this._arrow = $('<span class="flyout-arrow"></span>');
			this.element.append(this._arrow)
		},
        show: function(anchor, placement, alignment){
            this._baseFlyoutShow(anchor, placement, alignment)
        },
        //隐藏当前flyout对象
        hide: function(){
        	//一旦隐藏，就为document解除绑定的隐藏方法
        	this._documentBind(true);
			this._baseFlyoutHide();
        },
        lastAnchor: function( anchor ){
            return this._currentAnchor;
        },
        //创建停留计时器
        _createStayTimer: function(){
        	var self = this;
        	if(this.setting.stayTime){
        		this._clearStay();
        		this._stayTimer = setTimeout(function(){
        			self.hide()
        		}, this.setting.stayTime+this.setting.showDuration )
        	}
        },
        //清除计时器
        _clearStay: function(){
        	if(this._stayTimer){
    			clearTimeout( this._stayTimer );
    			this._stayTimer = null;
    		}
        },
		_documentName: '',
        _documentBind: function(unbind){
        	var self = this;

			if (unbind && this._documentName != '') {
				$document.unbind(this._documentName);
				self._documentName = '';
			} else if (this._documentName == '') {
				this._documentName = core.clickAnyWhereHideButMe(this.element, function () {
					self.hide();
				});
        	}
        },
        //删除flyout对象
        destroy: function(){
        	var id = this.primaryId;
        	this.element.remove();
        	Flyout.cache[ id ] = null;
        	delete Flyout.cache[ id ];
        },
        //恢复flyout的原始html，这将情况此操作前所有对flyout的修改
        recover: function(){
        	this.element.html( this._originalHTML )
        },
        _arrowPostion: function(arrow){
			var arrow = this.element.find('span.flyout-arrow').attr('class', 'flyout-arrow');
            var pos = '';
            switch( this._where.at ){
                case 'top':
                    pos = 'bottom'; break;
                case 'bottom':
                    pos = 'top'; break;
                case 'left':
                    pos = 'right'; break;
                case 'right':
                    pos = 'left'; break;
            }

			var arrowOffset = 15;
			if (this._currentAnchorBox.width > 15) {
				arrowOffset = 0;
			}

			if (this._currentAlignment == 'right') {
				this._startPosition.left += arrowOffset;
			}else if( this._currentAlignment == 'left' ){
                this._startPosition.left -= arrowOffset;
            }

            arrow.addClass(pos + ' ' + pos+'-'+this._currentAlignment);
        },
        _baseFlyoutShow: function( anchor, placement, alignment ){
        	
        	var self = this, element = this.element, sameAnchor = false;
        	
        	//动态更新锚点
        	if (!anchor) {
                anchor = this._anchor;
            }else if( !anchor.jquery ){
            	anchor = $( anchor );
            }
            //显示位置
            if (!placement) {
                placement = this._placement;
            }
            //对齐方式
            if (!alignment) {
                alignment = this._alignment;
            }
            
            if(!anchor){
				console.error('没有标记触发手柄');
            	return ;
            }else{
            	if( this._currentAnchor && this._currentAnchor[0] == anchor[0] ){
            		sameAnchor = true
            	}
            	this._currentAnchor = anchor;
            	this._currentPlacement = placement;
            	this._currentAlignment = alignment;
            }
            
            Flyout.currentShowFlyout = this;
            
            if(anchor){
            	//获取flyout显示位置
                this._getTopLeft();
            }
            
            if( !this.alone ){
            	this._hideAllVisibleFlyout();
            }

            //点击相同锚点是否隐藏
			if (sameAnchor && this.setting.sameAnchorHide && !this._isHide) {
				this.hide();
				//如果不是同一个锚点，或者同一个锚点却不可点，或者是同一个锚点可点但flyout是隐藏的
			} else if (!sameAnchor || ( sameAnchor && !this.setting.sameAnchorHide ) || ( sameAnchor && this.setting.sameAnchorHide && this._isHide )) {

            	//同一个flyout的不同手柄触发，要先移除之前的绑定事件， 但是要触发一次hide事件
                this._documentBind(true);

            	//自动隐藏，信息提示用
				this._createStayTimer();
				this.setting.arrow && self._arrowPostion();

				element.show().css('z-index', this._getZIndex(anchor)).css(this._startPosition).stop().animate(this._endPosition, {
            		easing: this.setting.easing,
                	duration: this.setting.showDuration,
                	complete: function(){
						self._isHide = false;
            			self._documentBind();
                	}
                });
            	
            	//当显示的时候调用
        		self.setting.onShow.call( self, element, anchor )
            }
        },
        _baseFlyoutHide: function(){

        	this._clearStay();
        	
        	var self = this, _hidePosition = { opacity: 0 }, hideOffset = this.setting.hideOffset;
			switch (this._where.at) {
				case 'bottom':
					_hidePosition.top = this.element.offset().top + hideOffset;
        			break;
				case 'top':
					_hidePosition.top = this.element.offset().top - hideOffset;
        			break;
        		case 'left':
					_hidePosition.left = this.element.offset().left - hideOffset;
        			break;
        		case 'right':
					_hidePosition.left = this.element.offset().left + hideOffset;
        			break;
        		default:
        			break;
        	}

			this.element.stop().animate(_hidePosition, {
				easing: this.setting.easing,
				duration: 160,
				complete: function(){
					$(this).hide();
					self._isHide = true;
					if(self.setting.destroy){
						self.destroy()
					}
				}
			});
        	//当隐藏时调用
        	this.setting.onHide.call( this, this.element, this._currentAnchor )
        },
        //隐藏除自己之外的所有flyout
        _hideAllVisibleFlyout: function(){
			var currentId = this.primaryId;
        	$.each( Flyout.cache, function(id, flyout){
				if (flyout && currentId != id && !flyout._isHide && !flyout.alone) {
        			flyout.hide();
        		}
        	})
        },
        //保持当前flyout一直可见
        _getZIndex: function(/* jQuery object */ anchor){
        	var zindex = 1000, anchorParentZindex = anchor.parentsUntil('body').last().css('z-index');
        	if( anchorParentZindex!= 'auto' && anchorParentZindex >= zindex ){
        		zindex = anchorParentZindex+1;
        	}
        	return zindex;
        },
        //获取位置坐标
        _getTopLeft: function(){
        	var currentAnchor = this._currentAnchor;
			var anchor = this._currentAnchorBox = {
        			left: currentAnchor.offset().left,
        			right: currentAnchor.offset().left + currentAnchor.outerWidth(),
        			top: currentAnchor.offset().top,
        			bottom: currentAnchor.offset().top + currentAnchor.outerHeight(),
        			width: currentAnchor.outerWidth(),
        			height: currentAnchor.outerHeight()
        		},
        		element = {
        			width: this.element.outerWidth(),
        			height: this.element.outerHeight()
        		};
        	
        	this._currentScrollTop = $window.scrollTop();
			this._startPosition = {opacity: 0};
        	this._endPosition = { opacity: 1 };
        	
        	switch(this._currentPlacement){
        		case 'top':
        			this._where.at = "top";
        			if(!this._fitTop( anchor, element )){
        				this._fitBottom( anchor, element );
        				this._where.at ="bottom";
        			}
        			this._alignHorizontally( anchor, element );
        			break;
        		case 'bottom':
        			this._where.at = "bottom";
        			if(!this._fitBottom( anchor, element )){
        				this._fitTop( anchor, element );
        				this._where.at = "top"
        			}
        			this._alignHorizontally( anchor, element );
        			break;
        		case "left":
        			this._where.at = "left";
                    if (!this._fitLeft(anchor, element)) {
                        // Didn't fit, needs scrollbar
                        this._fitRight(anchor, element);
                        this._where.at = "right"
                    }
                    this._alignVertically(anchor, element);
                    break;
                case "right":
                	this._where.at = "right";
                    if (!this._fitRight(anchor, element)) {
                         // Didn't fit, needs scrollbar
                        this._fitLeft(anchor, element);
                        this._where.at = "left";
                    }
                    this._alignVertically(anchor, element);
                    break;
        		case 'auto':
        			break;
        	}

            switch( this._currentAlignment ){
                case 'top':
                    break;
                case 'bottom':
                    break;
                case 'left':
                    break;
                case 'right':
                    break;
            }
        },
        //是否适合上面
        _fitTop: function(/* object */ anchorDimension,/* object */ flyoutDimension) {
            this._nextTop = anchorDimension.top - flyoutDimension.height;
            this._nextAnimTop = this._nextTop - this.setting.showOffset;
            return (this._nextTop >= 0 && this._nextAnimTop > this._currentScrollTop && this._nextAnimTop + flyoutDimension.height < anchorDimension.top );
        },
        //是否适合下面
        _fitBottom: function ( anchorDimension, flyoutDimension ) {
            this._nextTop = anchorDimension.bottom;
            this._nextAnimTop = anchorDimension.bottom + this.setting.showOffset;
            return ( this._nextTop >= 0 && this._nextAnimTop - this._currentScrollTop + flyoutDimension.height < $window.height() );
        },
        //是否适合左边
        _fitLeft: function( anchorDimension, flyoutDimension ){
        	this._nextLeft = anchorDimension.left - flyoutDimension.width + this.setting.showOffset;
            this._nextAnimLeft = this._nextLeft - this.setting.showOffset*2;
            return (this._nextLeft >= 0 && this._nextLeft + flyoutDimension.width < $window.width() );
        },
       //是否适合右边
        _fitRight: function (anchorDimension, flyoutDimension) {
            this._nextLeft = anchorDimension.right - this.setting.showOffset;
            this._nextAnimLeft = this._nextLeft + this.setting.showOffset*2;
            return ( this._nextLeft >= 0 && this._nextLeft + flyoutDimension.width < $window.width() );
        },
        //是否适合居左
        _fitAlignLeft: function(anchorDimension, flyoutDimension ){
            this._nextLeft = anchorDimension.left;
        	return anchorDimension.left + flyoutDimension.width < $window.width()
        },
        //是否适合居右
        _fitAlignRight: function(anchorDimension, flyoutDimension){
            this._nextLeft = anchorDimension.right - flyoutDimension.width;
        	return anchorDimension.right > flyoutDimension.width;
        },
        //纵向对齐
        _alignVertically: function( anchorDimension, flyoutDimension ){
            if( this._currentAlignment == 'center' ){
                this._nextTop = anchorDimension.top + anchorDimension.height / 2 - flyoutDimension.height / 2;
            }else if( this._currentAlignment == 'top' ){
                this._nextTop = anchorDimension.top
            }else if ( this._currentAlignment == 'bottom' ){
                this._nextTop = anchorDimension.bottom - flyoutDimension.height
            }

			this._startPosition.top = this._nextTop;
            this._startPosition.left = this._nextLeft;
			this._endPosition.left = this._nextAnimLeft;
        },
        //横向对齐方式
        _alignHorizontally: function( anchorDimension, flyoutDimension ){

            var width = $window.width();

            if( this._currentAlignment == 'center' ){
                this._nextLeft = anchorDimension.left + anchorDimension.width / 2 - flyoutDimension.width / 2;

                if( this._nextLeft + flyoutDimension.width > width ) {
                    this._currentAlignment = 'right';
                } else if( this._nextLeft < 0 ){
                    this._currentAlignment = 'left';
                }
            }

            if( this._currentAlignment == 'left' && !this._fitAlignLeft(anchorDimension, flyoutDimension) ){
                this._fitAlignRight(anchorDimension, flyoutDimension);
                this._currentAlignment = 'right';
            }else if( this._currentAlignment == 'right' && !this._fitAlignRight(anchorDimension, flyoutDimension) ){
                this._fitAlignLeft(anchorDimension, flyoutDimension);
                this._currentAlignment = 'left';
            }

        	this._startPosition.left = this._nextLeft;
			this._startPosition.top = this._nextTop;
			this._endPosition.top = this._nextAnimTop;
        }
    },{
    	primaryId: 0,
    	cache:{},
    	//只能聚焦一个flyout对象
    	currentShowFlyout: null,
    	//隐藏已存在的所有flyout
    	hideAllFlyout: function(){
    		$('.ui-flyout:visible').each(function(){
        		var $this = $(this), flyoutControl = $this.data('flyout');
        		if( flyoutControl && $this.is(':visible')){
            		flyoutControl.hide();
        		}
        	})
    	},
        template: {
            base : {
                element: '<div class="ui-flyout box"><div class="mod"><div class="bd"></div><div class="ft"></div></div></div>'
            }

        }
	});
	return Flyout
});