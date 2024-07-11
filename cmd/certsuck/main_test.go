package main

import "testing"

func Test_run(t *testing.T) {
	type args struct {
		opts *options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Jira",
			args: args{
				opts: &options{
					hostPort: "jira.belastingdienst.nl:443",
				},
			},
		},
		{
			name: "Jira der",
			args: args{
				opts: &options{
					hostPort: "jira.belastingdienst.nl:443",
					derOut: true,
				},
			},
		},
		{
			name: "Jira pretty options",
			args: args{
				opts: &options{
					showOpts: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_run_google(t *testing.T) {
	type args struct {
		opts *options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "insecure",
			args: args{
				opts: &options{
					hostPort: "mediaservernew.local:10443",
					insecure: true,
				},
			},
		},
		// {
		// 	name: "Jira pretty options",
		// 	args: args{
		// 		opts: &options{
		// 			hostPort: "www.google.com:443",
		// 			showOut: true,
		// 			noRoot: true,
		// 			noServer: true,
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
