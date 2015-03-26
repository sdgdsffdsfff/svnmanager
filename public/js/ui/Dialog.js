define([
'kernel',
'react',
'components/template/Dialog',
'bootstrap'
],
function( core, React, DialogComponent ){

    var Dialog = core.Class.extend({
        init: function( options, props ){

            var self = this;

            if (typeof props !== 'undefined') {
                $.extend(this, props);
            }

            this.options = $.extend({
                title: 'Dialog',
                showClass: 'fade',
                classStyle: '',
                size: '',
                onShow: $.noop,
                onShown: $.noop,
                onHide: $.noop,
                onHidden: $.noop,
                onRendered: $.noop,
                hideClass: ''
            }, options);

            this.element = $('<div class="modal" />')
                .addClass(this.options.showClass)
                .addClass(this.options.classStyle);

            this._reactElement = React.render( React.createElement(DialogComponent, {
                title: this.options.title
            }), this.element[0]);

            this.element.find('.modal-dialog').addClass(this.options.size);

            this.header = this.element.find('.modal-header');
            this.title = this.header.find('.modal-title');
            this.body = this.element.find('.modal-body');
            this.footer = this.element.find('.modal-footer');

            core.body.append( this.element );

            this.element.on('shown.bs.modal', function(){
                self.options.onShown.apply(self, arguments);
            });
            this.element.on('hidden.bs.modal', function(){
                self.options.onHidden.apply(self, arguments);
            });

            this.options.onRendered.call(self)
        },
        _reactElement: null,
        getReact: function(){
            return this._reactElement;
        },
        getRef: function( name ){

            if( name in this ){
                return this[name];
            }
            if( name in this._reactElement.refs ){
                return this._reactElement.refs[name];
            }
            if( typeof name === 'undefined' ){
                return this._reactElement.refs
            }
            return null;
        },
        setTitle: function( title ){
            this.title.html( title );
            return this;
        },
        show: function(){
            this.options.onShow.apply(this, arguments);
            this.element.modal('show');
            return this
        },
        hide: function(){
            this.options.onHide.apply(this, arguments);
            this.element.modal('hide');
            return this
        }
    });

    return Dialog;
});