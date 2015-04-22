define([
'kernel',
'react'
],
function( core, React ){

    var FormButton = React.createClass({
        getDefaultProps: function () {
            return {
                text: '',
                loadingText: 'Loading..',
                className: 'btn-default',
                click: $.noop
            }
        },
        disable: function () {
            this.$elem().addClass('disable')
        },
        enable: function(){
            this.$elem().removeClass('disable');
        },
        loading: function(text){
            if ( text ){
                this.$elem().data('loading-text', text);
            }
            this.$elem().button('loading');
        },
        reset: function(){
            this.$elem().button('reset');
        },
        setText: function( text ){
            this.$elem().html(text);
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
                <button
                    className={'btn btn-sm ' + this.props.className}
                    data-loading-text={this.props.loadingText}
                    onClick={this.props.click.bind(this.props.overload, this)}
                >{this.props.text}</button>
            )
        }
    });

    return FormButton;
});