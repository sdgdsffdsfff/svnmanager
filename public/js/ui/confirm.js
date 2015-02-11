define([
'kernel',
'ui/Flyout',
'react',
'components/ui/Confirm'
],
function(core, Flyout, React, Confirm){
    var flyout = null, text = null;

    return function(anchor, ok, options){

        options = $.extend({
            okText: 'Confirm',
            okClass: 'btn-primary',
            cancelText: 'Cancel',
            cancelClass: 'btn-default',
            html: 'Are you sure?',
            okClick: ok || $.noop
        }, options);

        if(flyout == null){

            var elem = $('<div class="ui-flyout confirm" />');
            React.render(React.createElement(Confirm, options), elem[0]);

            flyout = new Flyout(elem, {
                onRendered: function(element){
                    var thisFlyout = this;
                    element.delegate('button[role]','click',function(){
                        thisFlyout.hide();
                        var type = $(this).attr('role'), anchor = thisFlyout._currentAnchor;
                        if(type == 'ok'){
                            options.okClick.call( thisFlyout, anchor )
                        }
                        return false
                    });
                }
            });
            text = flyout.element.find('p.text');

            core.body.append( elem );
        }

        if( options.html ){
            text.css('width','auto').show().html( options.html );
            if(flyout.element.outerWidth() > 300 ){
                text.width(250);
            }
        }else{
            text.hide().html('');
        }

        flyout.show( anchor, "top");

        return flyout
    }
});