#!/usr/bin/env perl

use strict;
use warnings;
use HTTP::Daemon;
use HTTP::Status;
use Getopt::Long;

# Parse command line arguments
my $listen_addr = 'localhost:8080';
GetOptions('listen=s' => \$listen_addr) or die "Usage: $0 [--listen address:port]\n";

# Parse listen address
my ($host, $port) = ('', 8080);
if ($listen_addr =~ /^:(\d+)$/) {
    $port = $1;
} elsif ($listen_addr =~ /^(.+):(\d+)$/) {
    $host = $1;
    $port = $2;
} elsif ($listen_addr =~ /^(\d+)$/) {
    $port = $1;
}

# Default to localhost if not specified
$host = 'localhost' if $host eq '';

# Get client IP from headers
sub get_client_ip {
    my ($req) = @_;

    # Try X-Forwarded-For first
    my $forwarded = $req->header('X-Forwarded-For');
    if ($forwarded) {
        my @ips = split /,/, $forwarded;
        if (@ips) {
            $ips[0] =~ s/^\s+|\s+$//g;
            return $ips[0];
        }
    }

    # Try X-Real-IP
    my $real_ip = $req->header('X-Real-IP');
    return $real_ip if $real_ip;

    # Fall back to peer address
    return 'unknown';
}

# Create HTTP daemon
my $daemon = HTTP::Daemon->new(
    LocalAddr => $host,
    LocalPort => $port,
    ReuseAddr => 1,
) or die "Cannot create HTTP daemon: $!";

print STDERR "Starting server on $host:$port\n";

# Handle Ctrl+C gracefully
$SIG{INT} = sub {
    print STDERR "\nShutting down server...\n";
    exit 0;
};

# Main server loop
while (my $client = $daemon->accept) {
    while (my $req = $client->get_request) {
        # Check for Tailscale headers
        my $login_name = $req->header('Tailscale-User-Login');
        my $display_name = $req->header('Tailscale-User-Name');
        my $profile_pic_url = $req->header('Tailscale-User-Profile-Pic');

        # Create response
        my $res = HTTP::Response->new(HTTP::Status::RC_OK);
        $res->header('Cache-Control' => 'no-cache');

        if ($login_name && $display_name && $profile_pic_url) {
            # Serve HTML page with user information
            $res->header('Content-Type' => 'text/html');
            my $body = <<"HTML";
<!DOCTYPE html>
<html>
<head>
    <title>Hello from Perl!</title>
</head>
<body>
    <h1>(Perl) Tailscale User Information</h1>
    <p><strong>Login Name:</strong> $login_name</p>
    <p><strong>Name:</strong> $display_name</p>
    <p><strong>Profile Picture:</strong></p>
    <img src="$profile_pic_url" alt="Profile Picture" style="max-width: 200px;">
</body>
</html>
HTML
            $res->content($body);
        } else {
            # Print anonymous message with IP
            my $client_ip = get_client_ip($req);
            $res->header('Content-Type' => 'text/plain');
            $res->content("Hello anonymous from $client_ip\n");
        }

        # Send response
        $client->send_response($res);
    }

    $client->close;
    undef $client;
}
