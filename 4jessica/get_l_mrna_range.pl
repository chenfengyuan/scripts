#!/usr/bin/perl

use strict;
use warnings;
use Carp;
use feature qw{say};

# 一个参数:文件名
# 返回文件内容，去掉所有空白符
sub get_content{
  my $file=shift;
  open my $in,"<",$file;
  $_=join "",<$in>;
  s/\s//g;
  close $in;
  return $_;
}

# 一个参数:dna文件名
# 返回一个数组，每个元素也是一个数组，内容是dna exon的一段的开始也结束
sub get_dna_exon_range{
  my $dna_file=shift;
  my $dna_exon_content=&get_content($dna_file);
  my @dna_range=map {[ split /\.\./,$_ ]} split ",",$dna_exon_content;
  return @dna_range;
}

# 第一个参数:dna文件名
# 第二个参数:l_mrna文件名
# 返回一个数组，每个元素是a..b的一个字符串
sub get_l_mrna_range{
  my $dna_file=shift;
  my $l_mrna_file=shift;
  my @dna_range=get_dna_exon_range($dna_file);
  my $rna_length=get_content($l_mrna_file);
  my $total=0;
  my @intron_len;
  # 得到每个intron的长度
  for (my ($i,$len)=(0,scalar @dna_range);$i<$len - 1;$i++) {
    push @intron_len,($dna_range[$i+1]->[0]-$dna_range[$i]->[1]);
  }
  # reverse，这样从后面开始处理
  my @dna_range_reverse=reverse @dna_range;
  # 获取每个段的长度，存到@rna_len_reverse里面，这里会根据rna_length来获得长度，所以，要逆向处理。
  my @rna_len_reverse;
  foreach (@dna_range_reverse) {
    last if $total > $rna_length;
    my $len =$_->[1]-$_->[0];
    $total+=$len;
    push @rna_len_reverse,$len;
  }
  # 再次反转,获得正向的。然后通过长度推算出范围
  my @rna_len=reverse @rna_len_reverse;
  my @rna_range;
  my $i=1;
  foreach (@rna_len) {
    push @rna_range,"$i..@{[($i+$_)]}";
    $i+=$_+1;
  }
  return (\@rna_range,\@intron_len);
}

# 第一个参数:dna文件名
# 第二个参数:l_mrna文件名
# 返回结果字符串，调用get_l_mrna_range，然后转换成字符串返回
sub print_get_l_mrna_range{
  my $dna_file=shift;
  my $l_mrna_file=shift;
  my $output_file=shift;
  my @result=&get_l_mrna_range($dna_file,$l_mrna_file);
  my @rna_range=@{$result[0]};
  my @intron_len=@{$result[1]};
  open my $out,">",$output_file;
  say $out "exon number :";
  say $out scalar @rna_range;
  say $out "exon position :";
  say $out join ",",@rna_range;
  say $out "intron number :";
  say $out scalar @intron_len;
  say $out "intron size :";
  say $out join ",",@intron_len;
}
if(@ARGV<3){
  say "$0 : dna-file l_mrna_file output_file";
  exit(1);
}
my $dna_file=shift;
my $l_mrna_file=shift;
my $output_file=shift;
print_get_l_mrna_range($dna_file,$l_mrna_file,$output_file);

