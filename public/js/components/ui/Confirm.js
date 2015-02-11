define([
'kernel',
'react'
], function( core, React ){
    var Confirm = React.createClass({displayName: "Confirm",
        render: function(){
            return (
                React.createElement("div", null, 
                    React.createElement("p", {className: "text"}), 
                    React.createElement("div", {className: "buttons"}, 
                        React.createElement("button", {className: "btn btn-xs "+this.props.okClass, role: "ok"}, this.props.okText), 
                        React.createElement("button", {className: "btn btn-xs "+this.props.cancelClass, role: "cancel"}, this.props.cancelText)
                    )
                )
            )
        }
    });

    return Confirm;
});