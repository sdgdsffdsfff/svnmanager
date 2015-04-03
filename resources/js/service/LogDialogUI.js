/**
 * Created by languid on 3/26/15.
 */

define([
'kernel',
'angular',
'./module',
'ui/Dialog'
],
function (core, ng, service, Dialog) {
    service.factory('LogDialogUI', function () {
        var dialog = new Dialog({
            title: "Log",
            size: "modal-lg",
            onRendered: function(){
                var pre = $('<pre />');
                this.body.append( pre );
                this.pre = pre;
            }
        }, {
            setContent: function( html ){
                this.pre.html( html );
                return this
            }
        });
        return dialog;
    })
});
