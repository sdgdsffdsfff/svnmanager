define([
'kernel',
'react',
'ui/Dialog',
'components/form/FormBtns'
],
function( core, React, Dialog, FormBtns ){
    var UpfileList = React.createClass({
        getDefaultProps: function(){
            return {
                list: []
            }
        },
        render: function(){
            return (
                <div upfile-control>
                    <div className="control">
                        <span onClick={this.sortByAction} id="UpFileSortByAction" data-sortby="0">Action <i className="fa fa-sort"></i></span>
                        <span onClick={this.sortByPath} id="UpFileSortByPath" data-sortby="0">Path <i className="fa fa-sort"></i></span>
                        <span onClick={this.selectAll} data-all="false" id="UpFileSelectAllBtn">Select All</span>
                    </div>
                    <ul>
                        {this.props.list.map(function(item, index){
                            return (
                                <li key={index} data-id={item.Id}>
                                    <label>
                                        <span className={"action " + this.getAction(item.Action).toLowerCase()}>{this.getAction(item.Action)}</span>
                                        <span className="path">{item.Path}</span>
                                        <input type="checkbox" className="hidden" value={item.Id} onChange={this.check} />
                                        <span className="checkbox">
                                            <i className="fa fa-check"></i>
                                        </span>
                                    </label>
                                </li>
                            )
                        }, this)}
                    </ul>
                    <div className="info">
                        <span className="add">Add:<b>0</b></span>
                        <span className="update">Update:<b>0</b></span>
                        <span className="delete">Delete:<b>0</b></span>
                    </div>
                </div>
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
            this.$info = this.$el.find('.info');
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
        },
        selectAll: function(){
            var q = $.Deferred(), len = this.$checkboxes.length;
            if( !this.$selectAllBtn.data('all') ){
                this.$checkboxes.each(function(i){
                    core.delay(function(){
                        this.checked = true;
                        if( len-1 == i ){
                            q.resolve();
                        }
                    }.bind(this), i*16);
                });
            }else{
                this.$checkboxes.each(function(i){
                    core.delay(function(){
                        this.checked = false;
                        if( len-1 == i ){
                            q.resolve();
                        }
                    }.bind(this), i*16);
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
                q.reject()
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
            text: 'Deploy Now',
            className: 'btn-primary',
            click: events.confirm
        }];

        var dialog = new Dialog( $.extend(options, {
            title: 'Undeploy Files',
            classStyle: 'upgrade-dialog'
        }), extral);
        dialog.upfileList = React.render(<UpfileList />, dialog.body[0], function(){
            this.getView();
        });

        dialog.formBtns = React.render(<FormBtns
            buttons={buttons}
            overload={dialog}
        />, dialog.footer[0]);

        return dialog;
    }
});