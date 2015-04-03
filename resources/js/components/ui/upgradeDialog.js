define([
'kernel',
'react',
'ui/Dialog',
'components/form/FormBtns'
],
function( core, React, Dialog, FormBtns ){
    var UpfileList = React.createClass({displayName: "UpfileList",
        getDefaultProps: function(){
            return {
                list: []
            }
        },
        render: function(){
            return (
                React.createElement("div", null, 
                    React.createElement("div", {className: "control"}, 
                        React.createElement("span", {onClick: this.sortByAction, id: "UpFileSortByAction", "data-sortby": "0"}, React.createElement("i", {className: "fa fa-sort"}), " Action"), 
                        React.createElement("span", {onClick: this.sortByPath, id: "UpFileSortByPath", "data-sortby": "0"}, React.createElement("i", {className: "fa fa-sort"}), " Path"), 
                        React.createElement("span", {onClick: this.selectAll, "data-all": "false", id: "UpFileSelectAllBtn"}, "Select All")
                    ), 
                    React.createElement("ul", null, 
                        this.props.list.map(function(item, index){
                            return (
                                React.createElement("li", {key: index, "data-id": item.Id}, 
                                    React.createElement("label", null, 
                                        React.createElement("span", {className: "action " + this.getAction(item.Action).toLowerCase()}, this.getAction(item.Action)), 
                                        React.createElement("span", {className: "path"}, item.Path), 
                                        React.createElement("input", {type: "checkbox", className: "hidden", value: item.Id, onChange: this.check}), 
                                        React.createElement("span", {className: "checkbox"}, 
                                            React.createElement("i", {className: "fa fa-check"})
                                        )
                                    )
                                )
                            )
                        }, this)
                    ), 
                    React.createElement("div", {className: "info"}, 
                        React.createElement("span", {className: "add"}, "Add:", React.createElement("b", null, "0")), 
                        React.createElement("span", {className: "update"}, "Update:", React.createElement("b", null, "0")), 
                        React.createElement("span", {className: "delete"}, "Delete:", React.createElement("b", null, "0"))
                    ), 
                    React.createElement("div", {className: "message"}, 
                        React.createElement("textarea", {className: "form-control", placeholder: "Deploy Message"})
                    ), 
                    React.createElement("div", {className: "nofity"})
                )
            )
        },
        $el: null,
        $oldSort: null,
        $checkboxes: null,
        $items: null,
        redundancyMap: {},
        getView: function(){
            this.$el = $(this.getDOMNode());
            this.$list = this.$el.find('ul');
            this.$selectAllBtn = this.$el.find('#UpFileSelectAllBtn');
            this.$sortByActionBtn = this.$el.find('#UpFileSortByAction');
            this.$sortByPathBtn = this.$el.find('#UpFileSortByPath');
            this.$sortByVersionBtn = this.$el.find('#UpFileSortByVersion');
            this.$info = this.$el.find('.info');
            this.$nofity = this.$el.find('.nofity');
            this.$message = this.$el.find('textarea');
        },
        getCheckbox: function(){
            this.$items = this.$el.find('li');
            this.$checkboxes = this.$el.find('input:checkbox');
        },
        getAction: function( type ){
            return type == 1 ? 'Add' :
                   type == 2 ? 'Update' :
                   type == 3 ? 'Delete' : '';
        },
        getActionCount: function(list){
            var an = un = dn = 0;
            list.map(function( t ){
                if( t['Action'] == 1 ){
                    an++;
                }else if( t['Action'] == 2 ){
                    un++;
                }else if( t['Action'] == 3 ){
                    dn++;
                }
            });
            this.$info.find('.add b').html(an);
            this.$info.find('.update b').html(un);
            this.$info.find('.delete b').html(dn);
        },
        check: function(){
            if( this.$checkboxes.filter(':checked').length == this.$checkboxes.length ){
                this.$selectAllBtn.data('all', true).html('Deselect All');
            } else if( this.$selectAllBtn.data('all') ) {
                this.$selectAllBtn.data('all', false).html('Select All');
            }
        },
        message: function(){
            var q = $.Deferred(), value = this.$message.val().trim();
            if( value.length ) {
                q.resolve(value)
            }else{
                this.$message.focus();
                q.reject({
                    message: 'please enter deploy message'
                });
            }
            return q;
        },
        notify: function( text ){
            this.$nofity.html(text);
        },
        setList: function( list ){
            this.setProps({
                list: list
            });
            list.map(function(t){
                this.redundancyMap[t.Id] = t;
            }.bind(this));
            this.forceUpdate();
            this.$oldSort = this.$list.html();
            this.getCheckbox();
            this.getActionCount(list);
            this.sortBy('action', 1);
            this.selectAll( false );
        },
        selectAll: function( deselect ){
            var q = $.Deferred(), len = this.$checkboxes.length;
            if( !!deselect && !this.$selectAllBtn.data('all') ){
                this.$checkboxes.each(function(i){
                    core.delay(function(){
                        this.checked = true;
                        if( len-1 == i ){
                            q.resolve();
                        }
                    }.bind(this), i*30);
                });
            }else{
                this.$checkboxes.each(function(i){
                    core.delay(function(){
                        this.checked = false;
                        if( len-1 == i ){
                            q.resolve();
                        }
                    }.bind(this), i*30);
                });
            }
            q.then(function(){
                this.check();
            }.bind(this))
        },
        sortByAction: function(){
            this.sortBy('action');
        },
        sortByPath: function(){
            this.sortBy('path');
        },
        sortBy: function( type, by ){
            var anchor, sortby, icon;

            if( type == 'action' ){
                this.$sortByPathBtn.data('sortby', 0).find('i').attr('class', 'fa fa-sort');
                anchor = this.$sortByActionBtn;
            }else{
                this.$sortByActionBtn.data('sortby', 0).find('i').attr('class', 'fa fa-sort');
                anchor = this.$sortByPathBtn;
            }

            icon = anchor.find('i');
            sortby = by === undefined ? anchor.data('sortby') : by > 0 ? by - 1 : 0;

            if( sortby == 2 ){
                anchor.data('sortby', 0);
                this.$list.html( this.$oldSort );
                icon.attr('class', 'fa fa-sort');
                return;
            }

            var items = $.makeArray(this.$items);

            if( sortby == 0 ){
                items.sort(function(a, b){
                    var ta = $(a).find('span.'+type).text().toLowerCase(),
                        tb = $(b).find('span.'+type).text().toLowerCase();
                    return ta.localeCompare(tb);
                });
                anchor.data('sortby', 1);
                icon.attr('class','fa fa-sort-asc');
            }else if( sortby == 1) {
                items.sort(function(a, b){
                    var ta = $(a).find('span.'+type).text().toLowerCase(),
                        tb = $(b).find('span.'+type).text().toLowerCase();
                    return tb.localeCompare(ta);
                });
                anchor.data('sortby', 2);
                icon.attr('class','fa fa-sort-desc');
            }

            this.$list.html(items);
            this.getCheckbox();
        },
        getReadyToDeployFile: function(){
            var self = this;

            var q = $.Deferred(),
                checkedBoxes = this.$checkboxes.filter(':checked');

            if( checkedBoxes.length == 0 ){
                q.reject({
                    message: 'No file selected!'
                })
            } else if( checkedBoxes.length == this.$checkboxes.length) {
                //表示部署全部文件
                q.resolve([0])
            } else {
                var result = [];
                checkedBoxes.each(function(){
                    result.push( self.redundancyMap[this.value] )
                });
                q.resolve(result);
            }
            return q;
        }
    });

    return function( events, options, extral ){
        var buttons = [{
            text: 'Next',
            className: 'btn-primary',
            click: events.confirm
        }];

        var dialog = new Dialog( $.extend(options, {
            title: 'Undeploy Files',
            classStyle: 'upgrade-dialog'
        }), extral);
        dialog.upfileList = React.render(React.createElement(UpfileList, null), dialog.body[0], function(){
            this.getView();
        });

        dialog.formBtns = React.render(React.createElement(FormBtns, {
            buttons: buttons, 
            overload: dialog}
        ), dialog.footer[0]);

        return dialog;
    }
});