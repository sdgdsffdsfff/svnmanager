define([
'kernel',
'react'
], function( core, React ){
    var Dialog = React.createClass({displayName: "Dialog",
        render: function(){
            return (
                React.createElement("div", {className: "modal-dialog"}, 
                    React.createElement("div", {className: "modal-content"}, 
                        React.createElement("div", {className: "modal-header"}, 
                            React.createElement("button", {type: "button", className: "close", "data-dismiss": "modal", "aria-label": "Close"}, React.createElement("span", {"aria-hidden": "true"}, "×")), 
                            React.createElement("h4", {className: "modal-title"}, this.props.title)
                        ), 
                        React.createElement("div", {className: "modal-body"}), 
                        React.createElement("div", {className: "modal-footer"})
                    )
                )
            )
        }
    });

    return Dialog;
});