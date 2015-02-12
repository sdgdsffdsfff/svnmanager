/**
 * Created by languid on 12/26/14.
 */

define([
'kernel',
'react',
'ui/Flyout',
'components/form/helper',
'components/form/FormBody',
'components/form/FormBtns'
],
function (core, React, Flyout, formHelper, FormBody, FormBtns) {

    var FormFlyout = React.createClass({
        getDefaultProps: function () {
            return {
                fields: [],
                buttons: []
            }
        },
        render: function () {

            this.flyout = this.props.flyout;

            var titleEl = null;
            if( this.props.title ){
                titleEl = <div className="hd">{this.props.title}</div>
            }
            return (
                <div className="mod">
                    {titleEl}
                    <div className="bd">
                        <FormBody
                            fields={this.props.fields}
                            overload={this.props.overload}
                            inline={this.props.inline}
                            useLabel={this.props.useLabel}
                            ref="formBody"
                        />
                        <FormBtns
                            buttons={this.props.buttons}
                            overload={this.props.overload}
                            ref="formBtns"
                        />
                    </div>
                </div>
            )
        }
    });

    return function (config, flyoutConfig, extFlyout) {
        var btns = [{
            text: 'OK',
            className: 'btn-primary',
            click: function () {
                this.submitForm();
            }
        }];
        if (config) {
            if (config.buttons && config.buttons[0] == 'append') {
                config.buttons.splice(0, 1);
                config.buttons = btns.concat(config.buttons)
            } else {
                config.buttons = btns;
            }
        }
        config = $.extend({
            title: '',
            useLabel: false,
            inline: false,
            fields: [],
            buttons: []
        }, config);

        flyoutConfig = $.extend({
            onShow: $.noop,
            onHide: $.noop,
            init: $.noop,
            submit: $.noop
        }, flyoutConfig);

        extFlyout = $.extend({}, formHelper, extFlyout);

        var flyoutClass = 'form box';
        if (config.inline && config.useLabel) {
            flyoutClass += ' form-horizontal';
        }

        var div = $('<div />', {
            'class': flyoutClass
        });

        div.data('reactElement', React.render(<FormFlyout {...config} overload={flyout} />, div[0]));

        var flyout = new Flyout(div, flyoutConfig, extFlyout);

        flyout.arrow();

        document.body.appendChild(div[0]);

        return flyout;
    }
});