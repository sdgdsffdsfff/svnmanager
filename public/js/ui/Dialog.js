define([
'kernel',
'react',
'components/ui/Dialog',
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
                onShow: $.noop,
                onShown: $.noop,
                onHide: $.noop,
                onHidden: $.noop,
                hideClass: ''
            }, options);

            this.element = $('<div class="modal" />').addClass(this.options.showClass);

            this.dom = React.render( React.createElement(DialogComponent, {
                title: this.options.title
            }), this.element[0]);

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
        },
        setTitle: function( title ){
            this.title.html( title );
        },
        show: function(){
            this.options.onShow.apply(this, arguments);
            this.element.modal('show');
        },
        hide: function(){
            this.options.onHide.apply(this, arguments);
            this.element.modal('hide');
        }
    });

    return Dialog;
});