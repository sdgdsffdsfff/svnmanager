define([
'kernel',
'react',
'ui/Dialog',
'components/form/helper',
'components/form/FormBody',
'components/form/FormBtns',
'bootstrap'
],
function( core, React, Dialog, formHelper, FormBody, FormBtns ){

    return function( formConfig, dialogConfig, extDialog ){

        var btns = [{
            text: 'Confirm',
            className: 'btn-primary',
            click: function () {
                this.submitForm();
            }
        }];

        if (formConfig) {
            if (formConfig.buttons && formConfig.buttons[0] == 'append') {
                formConfig.buttons.splice(0, 1);
                formConfig.buttons = btns.concat(formConfig.buttons)
            } else {
                formConfig.buttons = btns;
            }
        }
        formConfig = $.extend({
            useLabel: false,
            inline: false,
            fields: [],
            buttons: []
        }, formConfig);

        dialogConfig = $.extend({
            onShow: $.noop,
            onHide: $.noop,
            init: $.noop,
            submit: $.noop
        }, dialogConfig);

        extDialog = $.extend({}, formHelper, extDialog);

        var dialog = new Dialog(dialogConfig, extDialog);

        dialog.formBody = React.render(React.createElement(FormBody, {
            fields: formConfig.fields, 
            overload: dialog, 
            inline: formConfig.inline, 
            useLabel: formConfig.useLabel}
        ), dialog.body[0]);

        dialog.formBtns = React.render(React.createElement(FormBtns, {
            buttons: formConfig.buttons, 
            overload: dialog}
        ), dialog.footer[0]);

        return dialog;
    }
});