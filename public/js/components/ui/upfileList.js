define(['kernel', 'react'], function( core, React ){
    var UpfileList = React.createClass({displayName: "UpfileList",
        getDefaultProps: function(){
            return {
                list: []
            }
        },
        setList: function( list ){
            this.setProps({
                list: list
            });
            this.forceUpdate();
        },
        getAction: function( type ){
            return type == 1 ? 'Add' :
                   type == 2 ? 'Update' :
                   type == 3 ? 'Delete' : '';
        },
        onclick: function( item ){
            console.log(item)
        },
        render: function(){
            return (
                React.createElement("div", {className: "mod upfile-list"}, 
                    React.createElement("div", {className: "hd"}, "Upgrade"), 
                    React.createElement("div", {className: "bd"}, 
                        React.createElement("ul", null, 
                            this.props.list.map(function(item, index){
                                return (
                                    React.createElement("li", {onClick: this.onclick.bind(this, item)}, 
                                        React.createElement("label", {className: "checkbox"}, 
                                            React.createElement("input", {type: "checkbox"}), 
                                            React.createElement("span", null, React.createElement("i", {className: "fa fa-check"}))
                                        ), 
                                        React.createElement("span", {className: "action " + this.getAction(item.Action).toLowerCase()}, this.getAction(item.Action)), 
                                        React.createElement("p", {className: "path"}, item.Path)
                                    )
                                )
                            }, this)
                        )
                    )
                )
            )
        }
    });

    return function(){
        var div = $('<div />', {
            id: 'id'+core.random(10)
        });
        return div.data('reactElement', React.render(React.createElement(UpfileList, null), div[0]));
    };
});