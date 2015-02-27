define([
'kernel',
'react'
],
function( core, React ){

    var FormButton = React.createClass({displayName: "FormButton",
        getDefaultProps: function () {
            return {
                text: '',
                className: 'btn-default',
                click: $.noop
            }
        },
        disable: function () {
            this.$elem().addClass('disable')
        },
        _$el: null,
        $elem: function(){
            if( !this._$el ){
                this._$el = $(this.getDOMNode());
            }
            return this._$el;
        },
        render: function () {
            return (
                React.createElement("button", {
                    className: 'btn btn-sm ' + this.props.className, 
                    onClick: this.props.click.bind(this.props.overload, this)
                }, this.props.text)
            )
        }
    });

    return FormButton;
});