define([
'kernel',
'react'
],
function( core, React ){

    var FormButton = React.createClass({
        getDefaultProps: function () {
            return {
                text: '',
                className: 'btn-default',
                click: $.noop
            }
        },
        disable: function () {
            $(this.getDOMNode()).addClass('disable')
        },
        render: function () {
            return (
                <button
                    className={'btn btn-sm ' + this.props.className}
                    onClick={this.props.click.bind(this.props.overload, this)}
                >{this.props.text}</button>
            )
        }
    });

    return FormButton;
});