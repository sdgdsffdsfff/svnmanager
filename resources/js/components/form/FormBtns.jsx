define([
'kernel',
'react',
'components/form/FormButton'
],
function( core, React, FormButton ){

    var FormBtns = React.createClass({
        getRef: function( name ){
            if( name in this.refs ){
                return this.refs[name]
            }
            if( typeof name === 'undefined' ){
                return this.refs;
            }
            return null;
        },
        render: function(){
            return (
                <div className="form-buttons">
                    {this.props.buttons.map(function (button, i) {
                        return <FormButton parent={this} overload={this.props.overload} ref={button.name} {...button} key={i} />;
                    }, this)}
                </div>
            )
        }
    });

    return FormBtns;
});